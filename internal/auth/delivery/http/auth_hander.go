package http

import (
	"log"
	"os"
	authService "simple-backend/internal/auth/service"
	model "simple-backend/internal/domain/auth"
	response "simple-backend/internal/utils/response"

	authUtils "simple-backend/internal/utils/auth"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type authController struct {
	service model.AuthServiceInterface
}

func AuthHandler(server *gin.RouterGroup, DB *gorm.DB) {
	controller := &authController{
		service: authService.Init(DB),
	}

	migrateUser(DB)

	server.POST("/login", controller.LoginHandler)
}

func (a *authController) LoginHandler(c *gin.Context) {
	var body *model.AuthTable

	if err := c.BindJSON(&body); err != nil {
		response.MakeErrorResponse(c, http.StatusUnprocessableEntity, err)
		return
	}

	user, err := a.service.GetUser(body)
	if err != nil {
		response.MakeErrorResponse(c, http.StatusNotFound, err)
		return
	}

	token, err := authUtils.GenerateToken(user.Name)
	if err != nil {
		response.MakeErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func migrateUser(DB *gorm.DB) {
	if DB.Migrator().HasTable(&model.AuthTable{}) {
		return
	}

	DB.AutoMigrate(&model.AuthTable{})

	pwdHashed := authUtils.GetPasswordHashed(os.Getenv("TEST_USER_PASSWORD"))

	var user = &model.AuthTable{}
	user.Name = os.Getenv("TEST_USER_NAME")
	user.Password = pwdHashed

	if err := DB.Create(user).Error; err != nil {
		log.Fatalln("Creat test user failed", err.Error())
	}
}
