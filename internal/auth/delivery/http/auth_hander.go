package http

import (
	"encoding/json"
	"log"
	"os"
	authService "simple-backend/internal/auth/service"
	authModel "simple-backend/internal/domain/auth"
	userModel "simple-backend/internal/domain/user"

	authUtils "simple-backend/internal/utils/auth"
	"simple-backend/internal/utils/databases"
	errorUtils "simple-backend/internal/utils/error"
	responseUtils "simple-backend/internal/utils/response"

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

// ShowAccount godoc
// @Summary      Login
// @Description  Login to get authrization jwt token
// @Tags         System
// @Accept       json
// @Produce      json
// @Param        data      body     authModel.LoginBody  true  "login info"
// @Success      200  {object}  response.Response
// @Router       /system/login [post]
func (a *authController) LoginHandler(c *gin.Context) {
	var body authModel.LoginBody
	customError := errorUtils.CheckErrAndConvert(c.BindJSON(&body), http.StatusUnprocessableEntity, nil, nil)

	if customError != nil {
		c.JSON(customError.HttpStatusCode, customError)
		return
	}

	token, customError := a.service.Login(&body)
	if customError != nil {
		c.JSON(customError.HttpStatusCode, customError)
		return
	}

	c.JSON(http.StatusOK, responseUtils.MakeCommonResponse(map[string]interface{}{"token": token}, nil, nil, nil))
}

// ShowAccount godoc
// @Summary      forgot-password
// @Description  Forgot password and reset it
// @Tags         System
// @Accept       json
// @Produce      json
// @Param        data   body   authModel.LoginBody  true  "forgot password info"
// @Success      200  {object}  response.Response
// @Router       /system/forgot-password [post]
func (a *authController) ForgotPasswordHandler(c *gin.Context) {
	data, err := c.GetRawData()

	customError := errorUtils.CheckErrAndConvert(err, http.StatusInternalServerError, nil, nil)
	if customError != nil {
		c.JSON(customError.HttpStatusCode, customError)
		return
	}

	var jsonData map[string]interface{}
	customError = errorUtils.CheckErrAndConvert(json.Unmarshal(data, &jsonData), http.StatusInternalServerError, nil, nil)
	if customError != nil {
		c.JSON(customError.HttpStatusCode, customError)
		return
	}

	// find user is exists
	customError = a.service.ForgotPassword(&authModel.LoginBody{
		Email: jsonData["email"].(string),
	})

	if customError != nil {
		c.JSON(customError.HttpStatusCode, customError)
		return
	}

	/*
		TODO:
		 1. return uri + token
		 2. add api endpoint, verify token, add redis cache and expired, return verify_code
		 3. reset password body struct should add verify_code
	*/
	responseCode := http.StatusAccepted
	c.JSON(responseCode, responseUtils.MakeCommonResponse(os.Getenv("RESET_PASSWORD_URI"), &responseCode, nil, nil))
}

func migrateUser(DB *gorm.DB) {
	DB.AutoMigrate(&userModel.UserTable{})

	user := &userModel.UserTable{
		Name:     os.Getenv("TEST_USER_NAME"),
		Email:    os.Getenv("TEST_USER_EMAIL"),
		Password: authUtils.GetPasswordHashed(os.Getenv("TEST_USER_PASSWORD")),
	}

	if err := DB.Create(user).Error; err != nil {
		log.Fatalln("Creat test user failed", err)
	}
}
