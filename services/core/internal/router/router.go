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
	"github.com/go-chi/cors"
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

	tokenService := services.NewJWTService([]byte(cfg.JWT.Secret), time.Duration(cfg.JWT.Expires)*time.Second)
	userService := services.NewUserService(userRepo, sessionRepo, tokenService)
	scheduleService := services.NewScheduleService(pool)

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
	subjectHandler := handlers.NewSubjectHandler(r.scheduleService)
	academicPeriodHandler := handlers.NewAcademicPeriodHandler(r.scheduleService)
	scheduleHandler := handlers.NewScheduleHandler(r.scheduleService)
	lessonLogHandler := handlers.NewLessonLogHandler(r.scheduleService)

	r.router.Use(middleware.Logger)
	r.router.Use(middleware.Recoverer)

	r.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.router.Post("/login", authHandler.Login)
	r.router.Post("/refresh", authHandler.RefreshToken)

	r.router.With(authMiddleware.ValidateToken).Group(func(r chi.Router) {
		r.Get("/logout", authHandler.Logout)

		r.Route("/admin", func(r chi.Router) {
			r.Use(authMiddleware.AdminRequire)

			r.Get("/sessions", authHandler.GetSessions)

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

		r.Route("/academic-periods", func(r chi.Router) {
			r.With(authMiddleware.AdminRequire).Group(func(r chi.Router) {
				r.Post("/", academicPeriodHandler.Create)
				r.Patch("/{id}", academicPeriodHandler.Update)
				r.Delete("/{id}", academicPeriodHandler.Delete)
			})
			r.Get("/{id}", academicPeriodHandler.GetByID)
			r.Get("/", academicPeriodHandler.GetAll)
		})

		r.Route("/groups", func(r chi.Router) {
			r.With(authMiddleware.AdminRequire).Group(func(r chi.Router) {
				r.Post("/", groupHandler.CreateGroup)
				r.Patch("/{id}", groupHandler.UpdateGroup)
				r.Delete("/{id}", groupHandler.DeleteGroup)
			})
			r.Get("/", groupHandler.GetAllGroups)
			r.Get("/{id}", groupHandler.GetGroup)
		})

		r.Route("/audiences", func(r chi.Router) {
			r.With(authMiddleware.AdminRequire).Group(func(r chi.Router) {
				r.Post("/", audienceHandler.CreateAudience)
				r.Patch("/{id}", audienceHandler.UpdateAudience)
				r.Delete("/{id}", audienceHandler.DeleteAudience)
			})
			r.Get("/{id}", audienceHandler.GetAudience)
			r.Get("/", audienceHandler.GetAllAudience)
		})

		r.Route("/subjects", func(r chi.Router) {
			r.With(authMiddleware.AdminRequire).Group(func(r chi.Router) {
				r.Post("/", subjectHandler.CreateSubject)
				r.Patch("/{id}", subjectHandler.UpdateSubject)
				r.Delete("/{id}", subjectHandler.DeleteSubject)
			})
			r.Get("/", subjectHandler.GetAllSubjects)
			r.Get("/group/{group_id}", subjectHandler.GetSubjectsByGroup)
			r.Get("/{id}", subjectHandler.GetSubject)
		})

		r.Route("/schedule", func(r chi.Router) {
			r.With(authMiddleware.AdminRequire).Group(func(r chi.Router) {
				r.Post("/", scheduleHandler.CreateScheduleTemplate)
				r.Patch("/{id}", scheduleHandler.UpdateScheduleTemplate)
				r.Delete("/{id}", scheduleHandler.DeleteScheduleTemplate)
				r.Post("/semester", scheduleHandler.CreateSemesterSchedule)
			})
			r.Get("/", scheduleHandler.GetAllScheduleTemplates)
			r.Get("/{id}", scheduleHandler.GetScheduleTemplate)
			r.Get("/group/{group_id}", scheduleHandler.GetGroupSchedule)
			r.Get("/teacher/{teacher_id}", scheduleHandler.GetTeacherSchedule)
			r.Get("/audience/{audience_id}", scheduleHandler.GetAudienceSchedule)
		})

		r.Route("/lessons", func(r chi.Router) {
			r.With(authMiddleware.AdminRequire).Group(func(r chi.Router) {
				r.Post("/cancel/{id}", lessonLogHandler.CancelLesson)
				r.Post("/reschedule", lessonLogHandler.RescheduleLesson)
				r.Post("/complete", lessonLogHandler.CompleteLessonFromDate)
			})
			r.Get("/", lessonLogHandler.GetLessonLogs)
		})
	})
}

func (r *Router) Router() *chi.Mux {
	return r.router
}
