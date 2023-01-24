package http

import (
	"errors"
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

	if !DB.Migrator().HasTable(&model.AuthTable{}) {
		migrateUser(DB)
	}

	server.POST("/system/login", controller.LoginHandler)
	server.POST("/system/create", controller.CreateHandler)
	server.POST("/system/update", controller.UpdateHandler)
	server.POST("/system/forgot-password", controller.ForgotPasswordHandler)
}

// Login
func (a *authController) LoginHandler(c *gin.Context) {
	var body model.AuthBody

	if err := c.BindJSON(&body); err != nil {
		response.MakeErrorResponse(c, http.StatusUnprocessableEntity, err)
		return
	}

	user, err := a.service.GetUser(&body)
	if err != nil {
		response.MakeErrorResponse(c, http.StatusNotFound, err)
		return
	}

	token, err := authUtils.GenerateToken(map[string]interface{}{
		"uid":      user.IdentityId,
		"userName": user.Name,
	})
	if err != nil {
		response.MakeErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Create User
func (a *authController) CreateHandler(c *gin.Context) {
	var body model.AuthBody
	if err := c.BindJSON(&body); err != nil {
		response.MakeErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if body.Password != body.ConfirmPassword {
		response.MakeErrorResponse(c, http.StatusBadRequest, errors.New("confirm_password doesn't match password"))
		return
	}

	user, err := a.service.CreateUser(&body)
	if err != nil {
		response.MakeErrorResponse(c, http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusCreated, response.MakeCommonResponse(user, http.StatusCreated))
}

// Update User
func (a *authController) UpdateHandler(c *gin.Context) {
	var body model.AuthBody
	if err := c.BindJSON(&body); err != nil {
		response.MakeErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	err := a.service.UpdateUser(&body)
	if err != nil {
		response.MakeErrorResponse(c, http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusAccepted, response.MakeCommonResponse("update success", http.StatusAccepted))
}

// Forgot Password
func (a *authController) ForgotPasswordHandler(c *gin.Context) {
	var body model.AuthBody
	if err := c.BindJSON(&body); err != nil {
		response.MakeErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	err := a.service.ForgotPassword(&body)
	if err != nil {
		response.MakeErrorResponse(c, http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusAccepted, response.MakeCommonResponse("please check your email", http.StatusAccepted))
}

func migrateUser(DB *gorm.DB) {
	DB.AutoMigrate(&model.AuthTable{})

	pwdHashed := authUtils.GetPasswordHashed(os.Getenv("TEST_USER_PASSWORD"))

	var user = &model.AuthTable{}
	user.Name = os.Getenv("TEST_USER_NAME")
	user.Password = pwdHashed

	if err := DB.Create(user).Error; err != nil {
		log.Fatalln("Creat test user failed", err.Error())
	}
}
