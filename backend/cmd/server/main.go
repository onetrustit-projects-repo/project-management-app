package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/pm-app/backend/internal/handlers"
	"github.com/pm-app/backend/internal/middleware"
	"github.com/pm-app/backend/internal/repositories"
	"github.com/pm-app/backend/internal/services"
	"github.com/pm-app/backend/internal/websocket"
)

func main() {
	// Load environment
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	db, err := repositories.NewDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := repositories.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize Redis
	redisClient := repositories.NewRedis(os.Getenv("REDIS_URL"))

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	workspaceRepo := repositories.NewWorkspaceRepository(db)
	projectRepo := repositories.NewProjectRepository(db)
	taskRepo := repositories.NewTaskRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, os.Getenv("JWT_SECRET"))
	workspaceService := services.NewWorkspaceService(workspaceRepo, userRepo)
	projectService := services.NewProjectService(projectRepo, workspaceRepo)
	taskService := services.NewTaskService(taskRepo, projectRepo, redisClient)
	dashboardService := services.NewDashboardService(taskRepo, projectRepo)

	// Initialize WebSocket hub
	wsHub := wshub.NewHub()
	go wsHub.Run()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	workspaceHandler := handlers.NewWorkspaceHandler(workspaceService)
	projectHandler := handlers.NewProjectHandler(projectService)
	taskHandler := handlers.NewTaskHandler(taskService, wsHub)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)
	wsHandler := handlers.NewWebSocketHandler(wsHub)

	// Setup router
	router := gin.Default()

	// Middleware
	router.Use(middleware.CORS())
	router.Use(middleware.RateLimiter(redisClient))

	// Health check
	router.GET("/health", handlers.HealthCheck)

	// API routes
	api := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/logout", authHandler.Logout)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.Auth(os.Getenv("JWT_SECRET")))
		{
			// User
			protected.GET("/users", handlers.GetUsers)
			protected.GET("/users/:id", handlers.GetUser)
			protected.PUT("/users/:id", handlers.UpdateUser)

			// Workspaces
			protected.GET("/workspaces", workspaceHandler.List)
			protected.POST("/workspaces", workspaceHandler.Create)
			protected.GET("/workspaces/:id", workspaceHandler.Get)
			protected.PUT("/workspaces/:id", workspaceHandler.Update)
			protected.DELETE("/workspaces/:id", workspaceHandler.Delete)
			protected.GET("/workspaces/:id/members", workspaceHandler.GetMembers)

			// Projects
			protected.GET("/projects", projectHandler.List)
			protected.POST("/projects", projectHandler.Create)
			protected.GET("/projects/:key", projectHandler.Get)
			protected.PUT("/projects/:key", projectHandler.Update)
			protected.DELETE("/projects/:key", projectHandler.Delete)
			protected.GET("/projects/:key/members", projectHandler.GetMembers)
			protected.POST("/projects/:key/members", projectHandler.AddMember)

			// Tasks
			protected.GET("/projects/:key/tasks", taskHandler.ListByProject)
			protected.POST("/projects/:key/tasks", taskHandler.Create)
			protected.GET("/tasks/:id", taskHandler.Get)
			protected.PUT("/tasks/:id", taskHandler.Update)
			protected.DELETE("/tasks/:id", taskHandler.Delete)
			protected.POST("/tasks/:id/move", taskHandler.Move)
			protected.POST("/tasks/:id/comments", taskHandler.AddComment)
			protected.GET("/tasks/:id/comments", taskHandler.GetComments)
			protected.GET("/tasks/:id/activity", taskHandler.GetActivity)

			// Dashboard
			protected.GET("/dashboard/stats", dashboardHandler.GetStats)
			protected.GET("/dashboard/projects/:key/burndown", dashboardHandler.GetBurndown)

			// WebSocket
			protected.GET("/ws", wsHandler.HandleWebSocket)
		}
	}

	// Start server
	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		log.Printf("Server starting on :%s", port)
		if err := router.Run(":" + port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}
