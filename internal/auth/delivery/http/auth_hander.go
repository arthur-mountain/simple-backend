package http

import (
	"encoding/json"
	"errors"
	"fmt"
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

func AuthHandler(server *gin.RouterGroup, db *databases.TMysql, REDIS *databases.MyRedis) {
	controller := &authController{
		service: authService.Init(db, REDIS),
	}

	db.Execute(func(DB *gorm.DB) error {
		if !DB.Migrator().HasTable(&userModel.UserTable{}) {
			migrateUser(DB)
		}
		return nil
	}, nil)

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
	customError := errorUtils.CheckErrAndConvert(
		c.BindJSON(&body),
		http.StatusUnprocessableEntity,
	)

	if customError != nil {
		c.JSON(customError.HttpStatusCode, customError)
		return
	}

	token, customError := a.service.Login(&body)
	if customError != nil {
		c.JSON(customError.HttpStatusCode, customError)
		return
	}

	responseUtils.New(map[string]interface{}{"token": token}).Done(c)
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

	customError := errorUtils.CheckErrAndConvert(err, http.StatusInternalServerError)
	if customError != nil {
		c.JSON(customError.HttpStatusCode, customError)
		return
	}

	var jsonData map[string]interface{}
	customError = errorUtils.CheckErrAndConvert(json.Unmarshal(data, &jsonData), http.StatusInternalServerError)
	if customError != nil {
		c.JSON(customError.HttpStatusCode, customError)
		return
	}

	email, ok := jsonData["email"]
	if !ok {
		customError = errorUtils.NewCustomError(
			errors.New("could found email field"),
			http.StatusBadRequest,
		)
		c.JSON(customError.HttpStatusCode, customError)
		return
	}
	// find user is exists
	customError = a.service.ForgotPassword(&authModel.LoginBody{
		Email: email.(string),
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
	responseUtils.New(os.Getenv("RESET_PASSWORD_URI")).SetHttpCode(http.StatusAccepted).Done(c)
}

func migrateUser(DB *gorm.DB) {
	user := new(userModel.UserTable)
	DB.AutoMigrate(user)

	user.Name = os.Getenv("TEST_USER_NAME")
	user.Email = os.Getenv("TEST_USER_EMAIL")
	user.Password = authUtils.GetPasswordHashed(os.Getenv("TEST_USER_PASSWORD"))

	if err := DB.Create(user).Error; err != nil {
		fmt.Println("Creat test user failed", err)
	}
}
