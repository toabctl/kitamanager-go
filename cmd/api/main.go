package main

import (
	"log"

	"github.com/eenemeene/kitamanager-go/internal/config"
	"github.com/eenemeene/kitamanager-go/internal/database"
	"github.com/eenemeene/kitamanager-go/internal/handlers"
	"github.com/eenemeene/kitamanager-go/internal/middleware"
	"github.com/eenemeene/kitamanager-go/internal/routes"
	"github.com/eenemeene/kitamanager-go/internal/store"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/eenemeene/kitamanager-go/docs"
)

// @title KitaManager API
// @version 1.0
// @description REST API for managing Users, Groups, and Organizations with RBAC support
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@kitamanager.example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	userStore := store.NewUserStore(db)
	groupStore := store.NewGroupStore(db)
	orgStore := store.NewOrganizationStore(db)
	employeeStore := store.NewEmployeeStore(db)
	childStore := store.NewChildStore(db)

	authHandler := handlers.NewAuthHandler(userStore, cfg.JWTSecret)
	userHandler := handlers.NewUserHandler(userStore)
	groupHandler := handlers.NewGroupHandler(groupStore)
	orgHandler := handlers.NewOrganizationHandler(orgStore)
	employeeHandler := handlers.NewEmployeeHandler(employeeStore)
	childHandler := handlers.NewChildHandler(childStore)

	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)

	r := gin.Default()

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.Setup(r, authHandler, userHandler, groupHandler, orgHandler, employeeHandler, childHandler, authMiddleware)

	log.Printf("Starting server on port %s", cfg.ServerPort)
	log.Printf("Swagger UI available at http://localhost:%s/swagger/index.html", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
