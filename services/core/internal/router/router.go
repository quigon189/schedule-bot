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
	groupRepo := repository.NewGroupRepo(pool)
	teacherRepo := repository.NewTeacherRepo(pool)
	studentRepo := repository.NewStudentRepo(pool)
	roleRepo := repository.NewRoleRepo(pool)

	tokenService := services.NewJWTService([]byte("123"), 24*time.Hour)
	userService := services.NewUserService(userRepo, sessionRepo, tokenService)
	scheduleService := services.NewScheduleService(audienceRepo, groupRepo, teacherRepo, studentRepo, userRepo, roleRepo)

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
	audienceHandler := handlers.NewAudienceHandler(r.scheduleService)
	groupHandler := handlers.NewGroupHandler(r.scheduleService)
	teacherHandler := handlers.NewTeacherHandler(r.scheduleService)
	studentHandler := handlers.NewStudentHandler(r.scheduleService)
	roleHandler := handlers.NewRoleHandler(r.scheduleService)

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
			r.Route("/roles", func(r chi.Router) {
				r.Post("/assign", roleHandler.AssignRole)
				r.Delete("/remove", roleHandler.RemoveRole)
			})

		})

		r.Route("/user", func(r chi.Router) {
			r.Get("/", userHandler.GetCurrentUser)
			r.Get("/{id}", userHandler.GetUser)
			r.Patch("/password", userHandler.UpdateUserPassword)
		})

		r.Route("/teachers", func(r chi.Router) {
			r.With(authMiddleware.AdminRequire).Group(func(r chi.Router) {
				r.Post("/", teacherHandler.CreateTeacher)
				r.Post("/assign", roleHandler.AssignTeacher)
				r.Delete("/{id}", teacherHandler.DeleteTeacher)
			})
			r.Get("/", teacherHandler.GetAllTeachers)
			r.Get("/{id}", teacherHandler.GetTeacher)
		})

		r.Route("/students", func(r chi.Router) {
			r.With(authMiddleware.AdminRequire).Group(func(r chi.Router) {
				r.Post("/", studentHandler.CreateStudent)
				r.Post("/assign", roleHandler.AssignStudent)
				r.Patch("/{id}/group", studentHandler.UpdateStudentGroup)
				r.Delete("/{id}", studentHandler.DeleteStudent)
			})
			r.Get("/", studentHandler.GetAllStudents)
			r.Get("/{id}", studentHandler.GetStudent)
		})

		r.Route("/groups", func(r chi.Router) {
			r.With(authMiddleware.AdminRequire).Group(func(r chi.Router) {
				r.Post("/", groupHandler.CreateGroup)
				r.Patch("/", groupHandler.UpdateGroup)
				r.Delete("/{id}", groupHandler.DeleteGroup)
			})
			r.Get("/", groupHandler.GetAllGroups)
			r.Get("/{id}", groupHandler.GetGroup)
		})

		r.Route("/audiences", func(r chi.Router) {
			r.With(authMiddleware.AdminRequire).Group(func(r chi.Router) {
				r.Post("/", audienceHandler.CreateAudience)
				r.Patch("/", audienceHandler.UpdateAudience)
				r.Delete("/{id}", audienceHandler.DeleteAudience)
			})
			r.Get("/{id}", audienceHandler.GetAudience)
			r.Get("/", audienceHandler.GetAllAudience)
		})

		r.With(authMiddleware.AdminRequire).Group(func(r chi.Router) {
		})
	})
}

func (r *Router) Router() *chi.Mux {
	return r.router
}
