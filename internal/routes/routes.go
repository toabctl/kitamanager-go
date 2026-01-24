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
	employeeHandler *handlers.EmployeeHandler,
	childHandler *handlers.ChildHandler,
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

			// Employees
			employees := protected.Group("/employees")
			{
				employees.GET("", employeeHandler.List)
				employees.GET("/:id", employeeHandler.Get)
				employees.POST("", employeeHandler.Create)
				employees.PUT("/:id", employeeHandler.Update)
				employees.DELETE("/:id", employeeHandler.Delete)
				employees.GET("/:id/contracts", employeeHandler.ListContracts)
				employees.GET("/:id/contracts/current", employeeHandler.GetCurrentContract)
				employees.POST("/:id/contracts", employeeHandler.CreateContract)
				employees.DELETE("/:id/contracts/:contractId", employeeHandler.DeleteContract)
			}

			// Children
			children := protected.Group("/children")
			{
				children.GET("", childHandler.List)
				children.GET("/:id", childHandler.Get)
				children.POST("", childHandler.Create)
				children.PUT("/:id", childHandler.Update)
				children.DELETE("/:id", childHandler.Delete)
				children.GET("/:id/contracts", childHandler.ListContracts)
				children.GET("/:id/contracts/current", childHandler.GetCurrentContract)
				children.POST("/:id/contracts", childHandler.CreateContract)
				children.DELETE("/:id/contracts/:contractId", childHandler.DeleteContract)
			}
		}
	}
}
