package router

import (
	"Clinic_backend/config"
	"Clinic_backend/internal/handler"
	"Clinic_backend/internal/middleware"
	"Clinic_backend/internal/repository"
	"Clinic_backend/internal/service"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(cfg *config.Config, db *pgxpool.Pool) *gin.Engine {
	if cfg.Env.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = cfg.Env.AllowedOrigins
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))

	r.Use(middleware.LoggerMiddleware())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Repos
	userRepo := repository.NewUserRepository(db)
	doctorRepo := repository.NewDoctorRepository(db)
	serviceRepo := repository.NewServiceRepository(db)
	serviceCategoryRepo := repository.NewServiceCategoryRepository(db)
	specRepo := repository.NewSpecializationRepository(db)
	scheduleRepo := repository.NewScheduleRepository(db)
	appointmentRepo := repository.NewAppointmentRepository(db)
	callbackRepo := repository.NewCallbackRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	auditRepo := repository.NewAuditRepository(db)

	// Services
	authService := service.NewAuthService(cfg, userRepo, refreshTokenRepo)
	doctorService := service.NewDoctorService(doctorRepo, specRepo, scheduleRepo)
	serviceService := service.NewServiceService(serviceRepo, serviceCategoryRepo, specRepo)
	serviceCategoryService := service.NewCategoryService(serviceCategoryRepo, specRepo)
	specializationService := service.NewSpecializationService(specRepo)
	scheduleService := service.NewScheduleService(scheduleRepo)
	emailService := service.NewEmailService(cfg)
	appointmentService := service.NewAppointmentService(appointmentRepo, userRepo, doctorRepo, scheduleRepo, emailService, cfg)
	callbackService := service.NewCallbackService(callbackRepo)
	auditService := service.NewAuditService(auditRepo)

	// 10 requests per minute per IP on auth
	authRateLimiter := middleware.RateLimiterMiddleware(10, time.Minute)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userRepo, auditService)
	doctorHandler := handler.NewDoctorHandler(doctorService, appointmentService, auditService)
	serviceHandler := handler.NewServiceHandler(serviceService)
	serviceCategoryHandler := handler.NewCategoryHandler(serviceCategoryService)
	specializationHandler := handler.NewSpecializationHandler(specializationService)
	scheduleHandler := handler.NewScheduleHandler(scheduleService)
	appointmentHandler := handler.NewAppointmentHandler(appointmentService)
	callbackHandler := handler.NewCallbackHandler(callbackService, auditService)

	api := r.Group("/api/v1")
	{
		// Auth routes (public) — rate limited
		auth := api.Group("/auth")
		auth.Use(authRateLimiter)
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/logout", authHandler.Logout)
		}

		// User routes
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware(cfg))
		{
			users.GET("/me", userHandler.GetMe)
			users.PUT("/me", userHandler.UpdateMe)
			users.GET("/:id", middleware.RoleMiddleware("admin", "doctor"), userHandler.GetByID)

			admin := users.Group("")
			admin.Use(middleware.RoleMiddleware("admin"))
			{
				admin.GET("", userHandler.GetAll)
				admin.PUT("/:id", userHandler.Update)
				admin.DELETE("/:id", userHandler.Delete)
			}
		}

		// Doctors routes
		doctors := api.Group("/doctors")
		{
			doctors.GET("/specialization/:id", doctorHandler.GetBySpecialization)
			doctors.GET("/:id/schedule", doctorHandler.GetDoctorSchedule)
			doctors.GET("/:id", doctorHandler.GetDoctorByID)
			doctors.GET("", doctorHandler.GetAllDoctors)

			doctorMe := doctors.Group("/me")
			doctorMe.Use(middleware.AuthMiddleware(cfg))
			doctorMe.Use(middleware.RoleMiddleware("doctor"))
			{
				doctorMe.GET("", doctorHandler.GetMyProfile)
				doctorMe.GET("/schedule", doctorHandler.GetMyScheduleHandler)
				doctorMe.GET("/appointments", doctorHandler.GetMyAppointmentsHandler)
			}

			doctorsAdmin := doctors.Group("")
			doctorsAdmin.Use(middleware.AuthMiddleware(cfg))
			doctorsAdmin.Use(middleware.RoleMiddleware("admin"))
			{
				doctorsAdmin.POST("", doctorHandler.CreateDoctor)
				doctorsAdmin.PUT("/:id", doctorHandler.UpdateDoctor)
				doctorsAdmin.DELETE("/:id", doctorHandler.DeleteDoctor)
			}
		}

		// Services routes
		services := api.Group("/services")
		{
			services.GET("/category/:id", serviceHandler.GetByCategory)
			services.GET("/specialization/:id", serviceHandler.GetBySpecialization)
			services.GET("/:id", serviceHandler.GetServiceByID)
			services.GET("", serviceHandler.GetAllServices)

			servicesAdmin := services.Group("")
			servicesAdmin.Use(middleware.AuthMiddleware(cfg))
			servicesAdmin.Use(middleware.RoleMiddleware("admin"))
			{
				servicesAdmin.POST("", serviceHandler.CreateService)
				servicesAdmin.PUT("/:id", serviceHandler.UpdateService)
				servicesAdmin.DELETE("/:id", serviceHandler.DeleteService)
			}
		}

		// Service Categories routes
		categories := api.Group("/service-categories")
		{
			categories.GET("", serviceCategoryHandler.GetAllCategories)
			categories.GET("/favorite", serviceCategoryHandler.GetFavorites)
			categories.GET("/:id", serviceCategoryHandler.GetCategoryByID)

			categoriesAdmin := categories.Group("")
			categoriesAdmin.Use(middleware.AuthMiddleware(cfg))
			categoriesAdmin.Use(middleware.RoleMiddleware("admin"))
			{
				categoriesAdmin.POST("", serviceCategoryHandler.CreateCategory)
				categoriesAdmin.PUT("/:id", serviceCategoryHandler.UpdateCategory)
				categoriesAdmin.PATCH("/:id/favorite", serviceCategoryHandler.ToggleFavorite)
				categoriesAdmin.DELETE("/:id", serviceCategoryHandler.DeleteCategory)
			}
		}

		// Specializations routes
		specializations := api.Group("/specializations")
		{
			specializations.GET("", specializationHandler.GetAllSpecializations)
			specializations.GET("/:id", specializationHandler.GetSpecializationByID)

			specializationsAdmin := specializations.Group("")
			specializationsAdmin.Use(middleware.AuthMiddleware(cfg))
			specializationsAdmin.Use(middleware.RoleMiddleware("admin"))
			{
				specializationsAdmin.POST("", specializationHandler.CreateSpecialization)
				specializationsAdmin.PUT("/:id", specializationHandler.UpdateSpecialization)
				specializationsAdmin.DELETE("/:id", specializationHandler.DeleteSpecialization)
			}
		}

		// Schedules routes
		schedules := api.Group("/schedules")
		{
			schedules.GET("/day/:day", scheduleHandler.GetByDay)
			schedules.GET("/:id", scheduleHandler.GetScheduleByID)

			schedulesAdmin := schedules.Group("")
			schedulesAdmin.Use(middleware.AuthMiddleware(cfg))
			schedulesAdmin.Use(middleware.RoleMiddleware("admin"))
			{
				schedulesAdmin.GET("", scheduleHandler.GetAllSchedules)
				schedulesAdmin.POST("", scheduleHandler.CreateSchedule)
				schedulesAdmin.PUT("/:id", scheduleHandler.UpdateSchedule)
				schedulesAdmin.DELETE("/:id", scheduleHandler.DeleteSchedule)
			}
		}

		// Callback requests routes
		callback := api.Group("/callback-requests")
		{
			callback.POST("", callbackHandler.CreateCallbackRequest)

			callbackAdmin := callback.Group("")
			callbackAdmin.Use(middleware.AuthMiddleware(cfg))
			callbackAdmin.Use(middleware.RoleMiddleware("admin"))
			{
				callbackAdmin.GET("", callbackHandler.GetAllCallbackRequests)
				callbackAdmin.GET("/:id", callbackHandler.GetCallbackRequestByID)
				callbackAdmin.PUT("/:id", callbackHandler.UpdateCallbackRequest)
				callbackAdmin.DELETE("/:id", callbackHandler.DeleteCallbackRequest)
			}
		}

		// Available slots (public)
		api.GET("/appointments/slots/:doctor_id", appointmentHandler.GetAvailableSlots)

		// Appointments routes (authenticated)
		appointments := api.Group("/appointments")
		appointments.Use(middleware.AuthMiddleware(cfg))
		{
			appointments.POST("", appointmentHandler.CreateAppointment)
			appointments.POST("/admin", appointmentHandler.AdminCreateAppointment)
			appointments.GET("/me", appointmentHandler.GetMyAppointments)
			appointments.PUT("/:id", appointmentHandler.UpdateAppointment)
			appointments.DELETE("/:id", appointmentHandler.CancelAppointment)
			appointments.POST("/:id/result", appointmentHandler.AddAppointmentResult)
			appointments.GET("/doctor/:doctor_id", appointmentHandler.GetDoctorAppointments)
			appointments.GET("", appointmentHandler.GetAllAppointments)
		}
	}

	return r
}
