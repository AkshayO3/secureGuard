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
	handler := &api.Handler{UserModel: userModel}

	r := gin.Default()

	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)
	r.POST("/logout", handler.Logout)
	r.POST("/refresh", handler.Refresh)

	err = r.Run(":" + os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Server startup failed: %v", err)
	}
}
