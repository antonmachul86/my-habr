package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"my-habr/services/auth/db"
	"my-habr/services/auth/handler"
	"my-habr/services/auth/repository"
	"my-habr/services/auth/service"
)

func main() {
	r := gin.Default()

	dsn := "postgres://user:pass@postgres:5432/habr?sslmode=disable"
	if err := db.InitDB(dsn); err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	userRepo := &repository.UserRepository{}
	authService := service.NewAuthService(userRepo, redisClient)
	authHandler := handler.NewAuthHandler(authService)

	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
	r.POST("/refresh", authHandler.Refresh)
	r.POST("/logout", authHandler.Logout)

	r.Run(":8001")
	fmt.Println("start")
}
