package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"toir-app/internal/config"
	"toir-app/internal/database"
	"toir-app/internal/handlers"
	"toir-app/internal/middleware"
	"toir-app/internal/models"
	"toir-app/internal/repository"
	"toir-app/internal/services"

	_ "toir-app/docs"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title ТОиР API
// @version 1.0
// @description API веб-приложения для управления техническим обслуживанием и ремонтом оборудования предприятия
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите JWT токен в формате: Bearer {token}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Equipment{},
		&models.RepairRequest{},
		&models.MaintenanceSchedule{},
		&models.MaintenanceLog{},
		&models.Part{},
		&models.WorkOrder{},
		&models.WorkOrderPart{},
	); err != nil {
		log.Fatalf("failed to auto-migrate: %v", err)
	}

	// Repositories
	userRepo := repository.NewGormUserRepository(db)
	equipmentRepo := repository.NewGormEquipmentRepository(db)
	repairRequestRepo := repository.NewGormRepairRequestRepository(db)
	scheduleRepo := repository.NewGormScheduleRepository(db)
	workOrderRepo := repository.NewGormWorkOrderRepository(db)

	// Services
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	equipmentService := services.NewEquipmentService(equipmentRepo)
	repairRequestService := services.NewRepairRequestService(repairRequestRepo)
	scheduleService := services.NewScheduleService(scheduleRepo)
	workOrderService := services.NewWorkOrderService(workOrderRepo)
	dashboardQuerier := services.NewGormDashboardQuerier(db)
	dashboardService := services.NewDashboardService(dashboardQuerier)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	equipmentHandler := handlers.NewEquipmentHandler(equipmentService)
	repairRequestHandler := handlers.NewRepairRequestHandler(repairRequestService)
	scheduleHandler := handlers.NewScheduleHandler(scheduleService)
	workOrderHandler := handlers.NewWorkOrderHandler(workOrderService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)
	userHandler := handlers.NewUserHandler(userRepo)

	// Echo
	e := echo.New()
	e.HideBanner = true

	e.Use(echomw.Logger())
	e.Use(echomw.Recover())
	e.Use(echomw.CORSWithConfig(echomw.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderAuthorization, echo.HeaderContentType},
	}))

	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Public routes
	e.POST("/api/auth/login", authHandler.Login)

	// Protected routes
	api := e.Group("/api", middleware.JWTAuth(cfg.JWTSecret))

	// Equipment
	api.GET("/equipment", equipmentHandler.List)
	api.POST("/equipment", equipmentHandler.Create, middleware.RequireRole("engineer", "admin"))
	api.GET("/equipment/:id", equipmentHandler.GetByID)
	api.PUT("/equipment/:id", equipmentHandler.Update, middleware.RequireRole("engineer", "admin"))
	api.DELETE("/equipment/:id", equipmentHandler.Delete, middleware.RequireRole("engineer", "admin"))

	// Repair Requests
	api.GET("/repair-requests", repairRequestHandler.List)
	api.POST("/repair-requests", repairRequestHandler.Create)
	api.GET("/repair-requests/:id", repairRequestHandler.GetByID)
	api.PUT("/repair-requests/:id", repairRequestHandler.Update)

	// Maintenance Schedules
	api.GET("/maintenance-schedules", scheduleHandler.List, middleware.RequireRole("engineer", "admin"))
	api.POST("/maintenance-schedules", scheduleHandler.Create, middleware.RequireRole("engineer", "admin"))
	api.PUT("/maintenance-schedules/:id", scheduleHandler.Update, middleware.RequireRole("engineer", "admin"))

	// Work Orders
	api.GET("/work-orders", workOrderHandler.List, middleware.RequireRole("engineer", "technician", "admin"))
	api.POST("/work-orders", workOrderHandler.Create, middleware.RequireRole("engineer", "technician", "admin"))
	api.PUT("/work-orders/:id", workOrderHandler.Update, middleware.RequireRole("engineer", "technician", "admin"))

	// Dashboard
	api.GET("/dashboard", dashboardHandler.GetMetrics, middleware.RequireRole("engineer", "admin"))

	// Users (admin only)
	api.GET("/users", userHandler.List, middleware.RequireRole("admin"))
	api.POST("/users", userHandler.Create, middleware.RequireRole("admin"))
	api.PUT("/users/:id", userHandler.Update, middleware.RequireRole("admin"))

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		addr := fmt.Sprintf(":%s", cfg.ServerPort)
		log.Printf("starting server on %s", addr)
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("failed to shutdown server: %v", err)
	}

	log.Println("server stopped")
}
