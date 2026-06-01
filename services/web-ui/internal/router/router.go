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
	alertsHandler := handlers.SetupAlertsHandler(5*time.Minute, 1*time.Minute)
	authHandler := handlers.NewAuthHandler(r.coreClient, r.store)
	dashboardHandler := handlers.NewDashboardHandler(r.coreClient)
	adminUsersHandler := handlers.NewAdminUsersHandler(r.coreClient, alertsHandler)
	groupsHandler := handlers.NewGroupsHandler(r.coreClient)
	subjectHandler := handlers.NewSubjectsHandler(r.coreClient)
	audiencesHandler := handlers.NewAudiencesHandler(r.coreClient)

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
		r.Get("/admin/users/upload-teachers-form", adminUsersHandler.UploadTeachersForm)
		r.Get("/admin/users/teacher-template", adminUsersHandler.DownloadTeacherTemplate)
		r.Post("/admin/users/upload-teachers", adminUsersHandler.UploadTeachers)

		r.Post("/admin/users", adminUsersHandler.CreateUser)

		r.Get("/admin/groups", groupsHandler.ListGroupsPage)
		r.Get("/admin/groups/list", groupsHandler.ListGroupsFragment)
		r.Get("/admin/groups/new", groupsHandler.NewGroupForm)
		r.Post("/admin/groups", groupsHandler.CreateGroup)
		r.Delete("/admin/groups/{id}", groupsHandler.DeleteGroup)
		r.Get("/admin/groups/upload-form", groupsHandler.UploadGroupsForm)
		r.Get("/admin/groups/template", groupsHandler.DownloadGroupTemplate)
		r.Post("/admin/groups/upload", groupsHandler.UploadGroupsExcel)

		r.Get("/admin/subjects", subjectHandler.SubjectsPage)
		r.Get("/admin/subjects/table", subjectHandler.SubjectsTable)

		r.Get("/admin/audiences", audiencesHandler.ListAudiencesPage)
		r.Get("/admin/audiences/table", audiencesHandler.TableFragment)
		r.Get("/admin/audiences/new", audiencesHandler.NewAudienceForm)
		r.Get("/admin/audiences/upload-form", audiencesHandler.UploadForm)
		r.Post("/admin/audiences", audiencesHandler.CreateAudience)
		r.Get("/admin/audiences/{id}/edit", audiencesHandler.EditAudienceForm)
		r.Post("/admin/audiences/{id}", audiencesHandler.UpdateAudience)
		r.Delete("/admin/audiences/{id}", audiencesHandler.DeleteAudience)
		r.Get("/admin/audiences/template", audiencesHandler.DownloadTemplate)
		r.Post("/admin/audiences/upload", audiencesHandler.UploadExcel)

	})

	r.router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Обработчик 404
	r.router.NotFound(handlers.NotFound)
}
