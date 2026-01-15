package router

import (
	"Clinic_backend/config"
	"Clinic_backend/internal/handler"
	"Clinic_backend/internal/middleware"
	"Clinic_backend/internal/repository"
	"Clinic_backend/internal/service"

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

	// CORS configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))

	// Logger middleware
	r.Use(middleware.LoggerMiddleware())

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Init Repos
	userRepo := repository.NewUserRepository(db)
	doctorRepo := repository.NewDoctorRepository(db)
	serviceRepo := repository.NewServiceRepository(db)
	serviceCategoryRepo := repository.NewServiceCategoryRepository(db)
	specRepo := repository.NewSpecializationRepository(db)
	scheduleRepo := repository.NewScheduleRepository(db)
	licenseRepo := repository.NewLicenseRepository(db)
	carouselRepo := repository.NewCarouselRepository(db)

	// Init Services
	authService := service.NewAuthService(cfg, userRepo)
	doctorService := service.NewDoctorService(doctorRepo, specRepo, scheduleRepo)
	serviceService := service.NewServiceService(serviceRepo, serviceCategoryRepo, specRepo)
	serviceCategoryService := service.NewCategoryService(serviceCategoryRepo, specRepo)
	specializationService := service.NewSpecializationService(specRepo)
	scheduleService := service.NewScheduleService(scheduleRepo)
	licenseService := service.NewLicenseService(licenseRepo)
	carouselService := service.NewCarouselService(carouselRepo)

	// Init handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userRepo)
	doctorHandler := handler.NewDoctorHandler(doctorService)
	serviceHandler := handler.NewServiceHandler(serviceService)
	serviceCategoryHandler := handler.NewCategoryHandler(serviceCategoryService)
	specializationHandler := handler.NewSpecializationHandler(specializationService)
	scheduleHandler := handler.NewScheduleHandler(scheduleService)
	licenseHandler := handler.NewLicenseHandler(licenseService)
	carouselHandler := handler.NewCarouselHandler(carouselService)

	api := r.Group("/api/v1")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// User routes
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware(cfg))
		{
			users.GET("/me", userHandler.GetMe)
			users.PUT("/me", userHandler.UpdateMe)

			// Admin only
			admin := users.Group("")
			admin.Use(middleware.RoleMiddleware("admin"))
			{
				admin.GET("", userHandler.GetAll)
				admin.GET("/:id", userHandler.GetByID)
				admin.PUT("/:id", userHandler.Update)
				admin.DELETE("/:id", userHandler.Delete)
			}
		}

		// Doctors routes
		doctors := api.Group("/doctors")
		{
			// Public routes
			doctors.GET("", doctorHandler.GetAllDoctors)
			doctors.GET("/:id", doctorHandler.GetDoctorByID)
			doctors.GET("/specialization/:id", doctorHandler.GetBySpecialization)
			doctors.GET("/:id/schedule", doctorHandler.GetDoctorSchedule)

			// Admin only
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
			// Public routes
			services.GET("", serviceHandler.GetAllServices)
			services.GET("/:id", serviceHandler.GetServiceByID)
			services.GET("/category/:id", serviceHandler.GetByCategory)
			services.GET("/specialization/:id", serviceHandler.GetBySpecialization)

			// Admin only
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
			// Public routes
			categories.GET("", serviceCategoryHandler.GetAllCategories)
			categories.GET("/:id", serviceCategoryHandler.GetCategoryByID)
			categories.GET("/favorite", serviceCategoryHandler.GetFavorites)

			// Admin only
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
			// Public routes
			specializations.GET("", specializationHandler.GetAllSpecializations)
			specializations.GET("/:id", specializationHandler.GetSpecializationByID)

			// Admin only
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
			// Public routes
			schedules.GET("/:id", scheduleHandler.GetScheduleByID)
			schedules.GET("/day/:day", scheduleHandler.GetByDay)

			// Admin only
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

		// Licenses routes
		licenses := api.Group("/licenses")
		{
			// Public routes
			licenses.GET("", licenseHandler.GetAllLicenses)
			licenses.GET("/:id", licenseHandler.GetLicenseByID)

			// Admin only
			licensesAdmin := licenses.Group("")
			licensesAdmin.Use(middleware.AuthMiddleware(cfg))
			licensesAdmin.Use(middleware.RoleMiddleware("admin"))
			{
				licensesAdmin.POST("", licenseHandler.CreateLicense)
				licensesAdmin.PUT("/:id", licenseHandler.UpdateLicense)
				licensesAdmin.DELETE("/:id", licenseHandler.DeleteLicense)
			}
		}

		// Carousel routes
		carousel := api.Group("/carousel")
		{
			// Public routes
			carousel.GET("", carouselHandler.GetAllSlides)
			carousel.GET("/:id", carouselHandler.GetSlideByID)

			// Admin only
			carouselAdmin := carousel.Group("")
			carouselAdmin.Use(middleware.AuthMiddleware(cfg))
			carouselAdmin.Use(middleware.RoleMiddleware("admin"))
			{
				carouselAdmin.POST("", carouselHandler.CreateSlide)
				carouselAdmin.PUT("/:id", carouselHandler.UpdateSlide)
				carouselAdmin.DELETE("/:id", carouselHandler.DeleteSlide)
			}
		}
	}

	return r
}
