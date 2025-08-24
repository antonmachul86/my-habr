package main

import (
	"my-habr/services/auth/db"
  "github.com/go-redis/redis/v9"
)

func main() {
	r := gin.Default()

	dsn := "postgres://user:pass@postgres:5432/habr?sslmode=disable"
    if err := db.InitBD(dsn); err!=nil{
      panic(err)
    }

    redisClient := redis.NewClient(&redis.Options){
      Addr: "redis:6379"
  }
}
