package http

import (
	"encoding/json"
	"log"
	"os"
	authService "simple-backend/internal/auth/service"
	authModel "simple-backend/internal/domain/auth"
	userModel "simple-backend/internal/domain/user"
	response "simple-backend/internal/utils/response"

	authUtils "simple-backend/internal/utils/auth"
	"simple-backend/internal/utils/databases"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type authController struct {
	service authModel.AuthServiceInterface
}

func AuthHandler(server *gin.RouterGroup, DB *gorm.DB, REDIS *databases.MyRedis) {
	controller := &authController{
		service: authService.Init(DB, REDIS),
	}

	if !DB.Migrator().HasTable(&userModel.UserTable{}) {
		migrateUser(DB)
	}

	server.POST("/system/login", controller.LoginHandler)
	server.POST("/system/forgot-password", controller.ForgotPasswordHandler)
}

// Login
func (a *authController) LoginHandler(c *gin.Context) {
	var body userModel.UserBody

	if err := c.BindJSON(&body); err != nil {
		response.MakeErrorResponse(c, http.StatusUnprocessableEntity, err)
		return
	}

	token, err := a.service.Login(&body)
	if err != nil {
		response.MakeErrorResponse(c, http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, response.MakeCommonResponse(struct{ token string }{token: token}))
}

// Forgot Password
func (a *authController) ForgotPasswordHandler(c *gin.Context) {
	data, err := c.GetRawData()
	if err != nil {
		response.MakeErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	var jsonData map[string]interface{}
	err = json.Unmarshal(data, &jsonData)

	if err != nil {
		response.MakeErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	err = a.service.ForgotPassword(&userModel.UserBody{
		Name: jsonData["name"].(string),
	})

	if err != nil {
		response.MakeErrorResponse(c, http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusAccepted, response.MakeCommonResponse(os.Getenv("RESET_PASSWORD_URI"), http.StatusAccepted))
}

func migrateUser(DB *gorm.DB) {
	DB.AutoMigrate(&userModel.UserTable{})

	pwdHashed := authUtils.GetPasswordHashed(os.Getenv("TEST_USER_PASSWORD"))

	var user = &userModel.UserTable{}
	user.Name = os.Getenv("TEST_USER_NAME")
	user.Password = pwdHashed

	if err := DB.Create(user).Error; err != nil {
		log.Fatalln("Creat test user failed", err.Error())
	}
}
