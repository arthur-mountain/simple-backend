package main

import (
	"fmt"
	"os"

	auth "simple-backend/internal/auth/delivery/http"
	authMiddleware "simple-backend/internal/auth/delivery/http/middleware"
	"simple-backend/internal/interactor/middleware/cors"
	"simple-backend/internal/interactor/middleware/statics"
	todo "simple-backend/internal/todo/delivery/http"
	"simple-backend/internal/utils/connect"

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

func main() {
	server := gin.Default()
	url := fmt.Sprintf(
		"%s:%s@%s(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_NETWORK"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)

	mysqlConfig := connect.MysqlConfig{DNS: &url}
	DB, err := mysqlConfig.Connect()
	if err != nil {
		fmt.Println("init db error", err.Error())
		return
	}

	// Apply middleware
	server.Use(cors.CorsMiddleware)
	server.Use(statics.StaticsMiddleware)

	serverV1 := server.Group("api/v1")
	// Login
	auth.AuthHandler(serverV1, DB)

	// Auth middleware
	serverV1.Use(authMiddleware.IsTokenValid)
	// Todo
	todo.TodoHandler(serverV1, DB)

	server.Run(fmt.Sprintf(":%s", os.Getenv("BACKEND_PORT")))
}
