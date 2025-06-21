package main

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"os/signal"
	"secureGuard/internal/api"
	"secureGuard/internal/data"
	"syscall"
	"time"
)

func main() {
	var err error
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Environment initialization failed -> %v", err.Error())
		return
	}
	db, errD := data.NewDB()
	if errD != nil {
		log.Fatalf("Cannot connect to the db %v", errD.Error())
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Couldn't close connection: %v", err.Error())
		}
	}(db)
	if errD := db.Ping(); errD != nil {
		log.Fatalf("DB Unreachable: %v", err)
	}
	userModel := &data.UserModel{DB: db}
	assetModel := &data.AssetModel{DB: db}
	vulnerabilityModel := &data.VulnerabilityModel{DB: db}
	incidentModel := &data.IncidentModel{DB: db}

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Next()
	})

	api.RegisterRoutes(router, userModel, assetModel, vulnerabilityModel, incidentModel)

	srv := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed %v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced shutdown:", err)
	}
	log.Println("Server exited")
}
