package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"sweng-task/internal/config"
	"sweng-task/internal/handler"
	"sweng-task/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	var (
		logger   *zap.Logger
		logLevel zapcore.Level
		err      error
	)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	if cfg.App.Environment == "development" {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		fmt.Printf("Error initializing logger: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		_ = logger.Sync()
	}()
	log := logger.Sugar()

	log.Infow("Configuration loaded",
		"environment", cfg.App.Environment,
		"log_level", cfg.App.LogLevel,
		"server_port", cfg.Server.Port,
	)

	// Parse log level from configuration
	if logLevel, err = zapcore.ParseLevel(cfg.App.LogLevel); err != nil {
		fmt.Printf("Error parsing log level: %v\n", err)
		os.Exit(1)
	}

	// Check log level increase
	if logLevel != log.Level() {
		log = log.WithOptions(zap.IncreaseLevel(logLevel))
	}

	// Initialize services
	lineItemService := service.NewLineItemService(log)
	adService := service.NewAdService(lineItemService, log)
	trackingService := service.NewTrackingService(lineItemService, log)

	// Setup Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Ad Bidding Service",
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout:  cfg.Server.Timeout,
	})

	// Register middleware
	app.Use(recover.New())
	app.Use(fiberlogger.New())
	app.Use(cors.New())

	// Register routes
	app.Get("/health", handler.HealthCheck)

	api := app.Group("/api/v1")

	// Line Item endpoints
	lineItemHandler := handler.NewLineItemHandler(lineItemService, log.Named("line_items"))
	api.Post("/lineitems", lineItemHandler.Create)
	api.Get("/lineitems", lineItemHandler.GetAll)
	api.Get("/lineitems/:id", lineItemHandler.GetByID).Name(handler.LineItemDetailsRoute)

	// Ad endpoints
	adHandler := handler.NewAdHandler(adService, log.Named("ads"))
	api.Get("/ads", adHandler.GetWinningAds)

	// Tracking endpoint
	trackingHandler := handler.NewTrackingHandler(trackingService, log.Named("tracking"))
	api.Post("/tracking", trackingHandler.TrackEvent)

	// Start server
	go func() {
		address := fmt.Sprintf(":%d", cfg.Server.Port)
		log.Infof("Starting server on %s", address)
		if err = app.Listen(address); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	if err = app.Shutdown(); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}

	log.Info("Server gracefully stopped")
}
