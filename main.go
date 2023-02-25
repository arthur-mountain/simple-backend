package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	auth "simple-backend/internal/auth/delivery/http"
	authMiddleware "simple-backend/internal/auth/delivery/http/middleware"
	"simple-backend/internal/interactor/middleware/cors"
	"simple-backend/internal/interactor/middleware/logs"
	"simple-backend/internal/interactor/middleware/statics"
	todo "simple-backend/internal/todo/delivery/http"
	user "simple-backend/internal/user/delivery/http"
	"simple-backend/internal/utils/databases"
	"simple-backend/internal/utils/swagger"

	"github.com/gin-gonic/gin"
)

// func GetPost(ctx *gin.Context) {
// 	id := ctx.Param("id")

// 	requestURL := fmt.Sprintf("https://jsonplaceholder.typicode.com/posts/%v", id)
// 	resp, err := http.Get(requestURL)
// 	if err != nil {
// 		log.Println("http.Get failed", err)
// 	}

// 	defer resp.Body.Close()

// 	body, _ := ioutil.ReadAll(resp.Body)
// 	fmt.Printf("data: %+v \n", string(body))

// 	ctx.Header("Content-Type", "application/json")
// 	ctx.String(http.StatusOK, string(body))
// }
// apiV1.GET("/todo/:id", GetPost)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	server := gin.Default()

	// mysql init
	DB_HANDLER := databases.MysqlInit(map[string]interface{}{
		"DNS": fmt.Sprintf(
			"%s:%s@%s(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			os.Getenv("MYSQL_USER"),
			os.Getenv("MYSQL_PASSWORD"),
			os.Getenv("MYSQL_NETWORK"),
			os.Getenv("MYSQL_HOST"),
			os.Getenv("MYSQL_PORT"),
			os.Getenv("MYSQL_DATABASE"),
		),
	})
	if err := DB_HANDLER.Connect(); err != nil {
		return
	}

	// redis init
	defaultRedisCtx := context.Background()
	defaultRedisDb, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	REDIS, err := databases.RedisInit(
		os.Getenv("REDIS_HOST")+":"+os.Getenv("REDIS_PORT"),
		"",
		defaultRedisDb,
		defaultRedisCtx,
	)
	if err != nil {
		fmt.Println("redis init error", err)
		return
	}

	// Apply middleware
	server.Use(cors.CorsMiddleware)
	server.Use(statics.StaticsMiddleware)
	server.Use(logs.LoggerToFile)

	serverV1 := server.Group("api/v1")
	// Login
	auth.AuthHandler(serverV1, DB_HANDLER, REDIS)

	// Auth middleware
	serverV1.Use(authMiddleware.IsTokenValid)
	// User
	user.UserHandler(serverV1, DB_HANDLER, REDIS)
	// Todo
	todo.TodoHandler(serverV1, DB_HANDLER)

	// Swagger set up
	swagger.Connect(server)

	server.Run(fmt.Sprintf(":%s", os.Getenv("BACKEND_PORT")))
}
