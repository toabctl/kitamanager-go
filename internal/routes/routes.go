package routes

import (
	"github.com/eenemeene/kitamanager-go/internal/handlers"
	"github.com/eenemeene/kitamanager-go/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Setup(
	r *gin.Engine,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	groupHandler *handlers.GroupHandler,
	orgHandler *handlers.OrganizationHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	api := r.Group("/api/v1")
	{
		api.POST("/login", authHandler.Login)

		protected := api.Group("")
		protected.Use(authMiddleware.RequireAuth())
		{
			// Organizations
			orgs := protected.Group("/organizations")
			{
				orgs.GET("", orgHandler.List)
				orgs.GET("/:id", orgHandler.Get)
				orgs.POST("", orgHandler.Create)
				orgs.PUT("/:id", orgHandler.Update)
				orgs.DELETE("/:id", orgHandler.Delete)
			}

			// Users
			users := protected.Group("/users")
			{
				users.GET("", userHandler.List)
				users.GET("/:id", userHandler.Get)
				users.POST("", userHandler.Create)
				users.PUT("/:id", userHandler.Update)
				users.DELETE("/:id", userHandler.Delete)
				users.POST("/:id/groups", userHandler.AddToGroup)
				users.DELETE("/:id/groups/:gid", userHandler.RemoveFromGroup)
				users.POST("/:id/organizations", userHandler.AddToOrganization)
				users.DELETE("/:id/organizations/:oid", userHandler.RemoveFromOrganization)
			}

			// Groups
			groups := protected.Group("/groups")
			{
				groups.GET("", groupHandler.List)
				groups.GET("/:id", groupHandler.Get)
				groups.POST("", groupHandler.Create)
				groups.PUT("/:id", groupHandler.Update)
				groups.DELETE("/:id", groupHandler.Delete)
				groups.POST("/:id/organizations", groupHandler.AddToOrganization)
				groups.DELETE("/:id/organizations/:oid", groupHandler.RemoveFromOrganization)
			}
		}
	}
}
