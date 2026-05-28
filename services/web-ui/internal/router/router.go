package router

import (
	"net/http"
	"time"
	"web-ui/internal/api"
	"web-ui/internal/handlers"
	"web-ui/internal/middlewares"
	"web-ui/internal/models"
	"web-ui/pkg/cache"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
)

type Router struct {
	coreClient *api.CoreClient
	router     *chi.Mux
	store      sessions.Store
	userCache  *cache.MemCache[models.User]
}

func NewRouter(client *api.CoreClient, store sessions.Store, userCache *cache.MemCache[models.User]) *Router {
	router := Router{
		coreClient: client,
		router:     chi.NewRouter(),
		store:      store,
		userCache:  userCache,
	}

	router.SetupRoutes()

	return &router
}

func (r *Router) Handler() *chi.Mux {
	return r.router
}

func (r *Router) SetupRoutes() {
	alertsHandler := handlers.NewAlertsHandler(5*time.Minute, 1*time.Minute)
	authHandler := handlers.NewAuthHandler(r.coreClient, r.store)
	dashboardHandler := handlers.NewDashboardHandler(r.coreClient)
	adminUsersHandler := handlers.NewAdminUsersHandler(r.coreClient, alertsHandler)
	groupsHandler := handlers.NewGroupsHandler(r.coreClient)
	subjectHandler := handlers.NewSubjectsHandler(r.coreClient)

	r.router.Use(middleware.Logger)
	r.router.Use(middlewares.RecoveryWithHTML)

	r.router.Get("/login", authHandler.LoginPage)
	r.router.Post("/login", authHandler.LoginSubmit)

	ms := middlewares.NewMiddlewares(r.store, r.coreClient, r.userCache)

	r.router.Group(func(r chi.Router) {
		//r.Use(csrfMiddleware)
		r.Use(ms.Auth)
		r.Get("/alerts", alertsHandler.GetAlerts)
		r.Get("/", dashboardHandler.Dashboard)
		r.Post("/logout", authHandler.Logout)
		r.Get("/admin/users", adminUsersHandler.ListUsersPage)
		r.Get("/admin/users/table", adminUsersHandler.TableFragment)
		r.Get("/admin/users/new", adminUsersHandler.NewUserForm)
		r.Post("/admin/users", adminUsersHandler.CreateUser)
		r.Get("/admin/groups", groupsHandler.ListGroupsPage)

		r.Get("/admin/subjects", subjectHandler.SubjectsPage)
		r.Get("/admin/subjects/table", subjectHandler.SubjectsTable)
	})

	r.router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Обработчик 404
	r.router.NotFound(handlers.NotFound)
}
