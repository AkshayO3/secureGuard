package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"secureGuard/internal/api"
	"secureGuard/internal/data"
	"secureGuard/internal/middleware"
)

func main() {
	var err error
	err = godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Environment initialization failed -> %v", err.Error())
		return
	}
	db, errD := sql.Open("postgres", os.Getenv("DB_URL"))
	if errD != nil {
		log.Panicf("Cannot connect to the db %v", errD.Error())
	}
	userModel := &data.UserModel{DB: db}
	assetModel := &data.AssetModel{DB: db}
	vulnerabilityModel := &data.VulnerabilityModel{DB: db}
	incidentModel := &data.IncidentModel{DB: db}
	handler1 := &api.UserHandler{UserModel: userModel}
	handler2 := &api.AssetHandler{AssetModel: assetModel}
	handler3 := &api.VulnerabilityHandler{VulnerabilityModel: vulnerabilityModel}
	handler4 := &api.IncidentHandler{IncidentModel: incidentModel}

	r := gin.Default()

	r.POST("/register", handler1.Register)
	r.POST("/login", handler1.Login)
	r.POST("/logout", middleware.RequireRole("viewer"), handler1.Logout)
	r.POST("/refresh", middleware.RequireRole("viewer"), handler1.Refresh)
	r.GET("/users", middleware.RequireRole("admin"), handler1.List)
	r.GET("/users/:id", middleware.RequireRole("admin"), handler1.GetById)
	r.PUT("/users/:id", middleware.RequireRole("admin"), handler1.UpdateById)
	r.DELETE("/users/:id", middleware.RequireRole("admin"), handler1.DeleteById)

	r.GET("/assets", handler2.ListAssets)
	r.GET("/assets/:id", handler2.GetById)
	r.POST("/assets", middleware.RequireRole("viewer"), handler2.Create)
	r.PUT("/assets/:id", middleware.RequireRole("admin"), handler2.Update)
	r.DELETE("/assets/:id", middleware.RequireRole("admin"), handler2.Delete)
	r.GET("/assets/:id/vulnerabilities", handler2.GetAssociatedV)
	r.GET("/assets/:id/incidents", handler2.GetAssociatedI)

	r.GET("/vulnerabilities", handler3.ListVulnerabilities)
	r.GET("/vulnerabilities/:id", handler3.Get)
	r.POST("/vulnerabilities", middleware.RequireRole("viewer"), handler3.Create)
	r.PUT("/vulnerabilities/:id", middleware.RequireRole("admin"), handler3.Update)
	r.DELETE("/vulnerabilities/:id", middleware.RequireRole("admin"), handler3.Delete)
	r.GET("/vulnerabilities/:id/assets", handler3.GetAssociatedAssets)
	r.POST("/vulnerabilities/assets", middleware.RequireRole("admin"), handler3.AssociateAssetWithVuln)

	r.GET("/incidents", handler4.ListIncidents)
	r.GET("/incidents/:incidentId", handler4.GetById)
	r.POST("/incidents", middleware.RequireRole("viewer"), handler4.CreateIncident)
	r.PUT("/incidents/:incidentId", middleware.RequireRole("admin"), handler4.UpdateById)
	r.DELETE("/incidents/:incidentId", middleware.RequireRole("admin"), handler4.DeleteById)
	r.GET("/incidents/:incidentId/assets", handler4.ListAssociatedAssets)
	r.POST("/incidents/assets", middleware.RequireRole("viewer"), handler4.AssociateIncidentWithAsset)

	err = r.Run(":" + os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Server startup failed: %v", err)
	}
}
