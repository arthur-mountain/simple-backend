package swagger

import (
	"fmt"
	"os"
	"simple-backend/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	docs.SwaggerInfo.Title = "Simple Backend API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Description = "This is a simple backend server."
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", os.Getenv("BACKEND_PORT"))
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}

func Connect(server *gin.Engine) {
	server.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
