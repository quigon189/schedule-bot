package router

import (
	"net/http"
	"web-ui/internal/api"
	"web-ui/internal/handlers"
	"web-ui/internal/middlewares"
	"web-ui/internal/session"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	sessionManager *session.SessionManager
	coreClient     *api.CoreClient
	router         *chi.Mux
}

func NewRouter(client *api.CoreClient, sm *session.SessionManager) *Router {
	router := Router{
		coreClient: client,
		sessionManager: sm,
		router: chi.NewRouter(),
	}

	router.SetupRoutes()

	return &router
}

func (r *Router) Handler() *chi.Mux {
	return r.router
}

func (r *Router) SetupRoutes() {
	sm := r.sessionManager
	authHandler := handlers.NewAuthHandler(r.coreClient, r.sessionManager)	
	dashboardHandler := handlers.NewDashboardHandler(r.coreClient, r.sessionManager)
	adminUsersHandler := handlers.NewAdminUsersHandler(r.coreClient, r.sessionManager)

	r.router.Use(middleware.Logger)
	r.router.Use(middleware.Recoverer)

	r.router.Get("/login", authHandler.LoginPage)
	r.router.Post("/login", authHandler.LoginSubmit)

	r.router.Group(func(r chi.Router) {
		//r.Use(csrfMiddleware)
		r.Use(middlewares.AuthMiddleware(sm))
		r.Get("/", dashboardHandler.Dashboard)
		r.Post("/logout", authHandler.Logout)
		r.Get("/admin/users", adminUsersHandler.ListUsersPage)
		r.Get("/admin/users/table", adminUsersHandler.TableFragment)
		r.Get("/admin/users/new", adminUsersHandler.NewUserForm)
		r.Post("/admin/users", adminUsersHandler.CreateUser)
	})

	r.router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
}
