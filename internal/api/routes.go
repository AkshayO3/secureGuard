package api

import (
	"github.com/gin-gonic/gin"
	"secureGuard/internal/data"
	"secureGuard/internal/middleware"
)

func RegisterRoutes(router *gin.Engine, userModel *data.UserModel, assetModel *data.AssetModel,
	vulnerabilityModel *data.VulnerabilityModel, incidentModel *data.IncidentModel) {

	userHandler := &UserHandler{UserModel: userModel}
	assetHandler := &AssetHandler{AssetModel: assetModel}
	vulnerabilityHandler := &VulnerabilityHandler{VulnerabilityModel: vulnerabilityModel}
	incidentHandler := &IncidentHandler{IncidentModel: incidentModel}

	userRoutes := router.Group("/users")
	{
		userRoutes.POST("/register", userHandler.Register)
		userRoutes.POST("/login", userHandler.Login)
		userRoutes.POST("/logout", middleware.RequireRole("viewer"), userHandler.Logout)
		userRoutes.POST("/refresh", middleware.RequireRole("viewer"), userHandler.Refresh)
		userRoutes.GET("", middleware.RequireRole("admin"), userHandler.List)
		userRoutes.GET("/:id", middleware.RequireRole("admin"), userHandler.GetById)
		userRoutes.PUT("/:id", middleware.RequireRole("admin"), userHandler.UpdateById)
		userRoutes.DELETE("/:id", middleware.RequireRole("admin"), userHandler.DeleteById)
	}

	assetRoutes := router.Group("/assets")
	{
		assetRoutes.GET("", assetHandler.List)
		assetRoutes.GET("/:id", assetHandler.Get)
		assetRoutes.POST("", middleware.RequireRole("viewer"), assetHandler.Insert)
		assetRoutes.PUT("/:id", middleware.RequireRole("admin"), assetHandler.Update)
		assetRoutes.DELETE("/:id", middleware.RequireRole("admin"), assetHandler.Delete)
		assetRoutes.GET("/:id/vulnerabilities", assetHandler.GetAssociatedV)
		assetRoutes.GET("/:id/incidents", assetHandler.GetAssociatedI)
	}

	vulnRoutes := router.Group("/vulnerabilities")
	{
		vulnRoutes.GET("", vulnerabilityHandler.List)
		vulnRoutes.GET("/:id", vulnerabilityHandler.Get)
		vulnRoutes.POST("", middleware.RequireRole("viewer"), vulnerabilityHandler.Insert)
		vulnRoutes.PUT("/:id", middleware.RequireRole("admin"), vulnerabilityHandler.Update)
		vulnRoutes.DELETE("/:id", middleware.RequireRole("admin"), vulnerabilityHandler.Delete)
		vulnRoutes.GET("/:id/assets", vulnerabilityHandler.GetAssociatedAssets)
		vulnRoutes.POST("/assets", middleware.RequireRole("admin"), vulnerabilityHandler.AssociateA)
	}

	incidentRoutes := router.Group("/incidents")
	{
		incidentRoutes.GET("", incidentHandler.List)
		incidentRoutes.GET("/:incidentId", incidentHandler.Get)
		incidentRoutes.POST("", middleware.RequireRole("viewer"), incidentHandler.Insert)
		incidentRoutes.PUT("/:incidentId", middleware.RequireRole("admin"), incidentHandler.Update)
		incidentRoutes.DELETE("/:incidentId", middleware.RequireRole("admin"), incidentHandler.Delete)
		incidentRoutes.GET("/:incidentId/assets", incidentHandler.ListAssociatedAssets)
		incidentRoutes.POST("/assets", middleware.RequireRole("viewer"), incidentHandler.AssociateA)
	}

}
