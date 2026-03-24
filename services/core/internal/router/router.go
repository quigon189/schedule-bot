package router

import (
	"core/internal/config"
	"core/internal/handlers"
	"core/internal/middlewares"
	"core/internal/repository"
	"core/internal/services"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Router struct {
	tokenService    *services.JWTService
	authService     *services.UserService
	scheduleService *services.ScheduleService
	router          *chi.Mux
}

func New(cfg *config.Config, pool *pgxpool.Pool) *Router {
	userRepo := repository.NewUserRepo(pool)
	sessionRepo := repository.NewSessionRepo(pool)
	audienceRepo := repository.NewAudienceRepo(pool)

	tokenService := services.NewJWTService([]byte("123"), 24*time.Hour)
	userService := services.NewUserService(userRepo, sessionRepo, tokenService)
	scheduleService := services.NewScheduleService(audienceRepo)

	router := Router{
		tokenService:    tokenService,
		authService:     userService,
		scheduleService: scheduleService,
		router:          chi.NewRouter(),
	}

	router.SetupRoutes()

	return &router
}

func (r *Router) SetupRoutes() {
	authMiddleware := middlewares.NewAuthMiddleware(r.authService)

	authHandler := handlers.NewAuthHandler(r.authService)
	userHandler := handlers.NewUserHandler(r.authService)
	scheduleHandler := handlers.NewScheduleHandler(r.scheduleService)

	r.router.Use(middleware.Logger)
	r.router.Use(middleware.Recoverer)

	r.router.Post("/login", authHandler.Login)
	r.router.Post("/refresh", authHandler.RefreshToken)

	r.router.With(authMiddleware.ValidateToken).Group(func(r chi.Router) {
		r.Get("/logout", authHandler.Logout)

		r.Route("/admin", func(r chi.Router) {
			r.Use(authMiddleware.AdminRequire)
			r.Route("/users", func(r chi.Router) {
				r.Get("/", userHandler.GetUsers)
				r.Get("/{id}", userHandler.GetUser)
				r.Post("/", userHandler.CreateUser)
				r.Patch("/password", userHandler.UpdateUserPassword)
			})
		})

		r.Route("/user", func(r chi.Router) {
			r.Get("/", userHandler.GetCurrentUser)
			r.Get("/{id}", userHandler.GetUser)
			r.Patch("/password", userHandler.UpdateUserPassword)
		})

		r.Route("/audiences", func(r chi.Router) {
			r.With(authMiddleware.AdminRequire).Group(func(r chi.Router) {
				r.Post("/", scheduleHandler.CreateAudience)
				r.Patch("/", scheduleHandler.UpdateAudience)
				r.Delete("/{id}", scheduleHandler.DeleteAudience)
			})
			r.Get("/{id}", scheduleHandler.GetAudience)
			r.Get("/", scheduleHandler.GetAllAudience)
		})

		r.With(authMiddleware.AdminRequire).Group(func(r chi.Router) {
		})
	})
}

func (r *Router) Router() *chi.Mux {
	return r.router
}
