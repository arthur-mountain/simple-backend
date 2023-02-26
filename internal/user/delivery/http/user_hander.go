package http

import (
	model "simple-backend/internal/domain/user"
	userService "simple-backend/internal/user/service"
	errorUtils "simple-backend/internal/utils/error"
	responseUtils "simple-backend/internal/utils/response"
	"strconv"

	"simple-backend/internal/utils/databases"

	"net/http"

	"github.com/gin-gonic/gin"
)

type userController struct {
	service model.UserServiceInterface
}

func UserHandler(server *gin.RouterGroup, db *databases.TMysql, REDIS *databases.MyRedis) {
	controller := &userController{
		service: userService.Init(db, REDIS),
	}

	// Get All User
	server.GET("/users", controller.GetUsers)
	// Get User
	server.GET("/users/:id", controller.GetUser)
	// Create User
	server.POST("/users/create", controller.CreateUser)
	// Update User
	server.PUT("/users/:id", controller.UpdateUser)
	// Delete User
	server.DELETE("/users/:id", controller.DeleteUser)
}

// ShowAccount godoc
// @Summary      Users
// @Description  All of user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /users [get]
func (a *userController) GetUsers(c *gin.Context) {
	users, customError := a.service.GetUsers()

	if customError != nil {
		c.JSON(customError.HttpStatusCode, customError)
		return
	}

	responseUtils.New(users).Done(c)
}

// ShowAccount godoc
// @Summary      User
// @Description  Get user by id
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path    int   true    "user id"
// @Success      200  {object}  response.Response
// @Router       /users/{id} [get]
func (a *userController) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	customError := errorUtils.CheckErrAndConvert(
		err,
		http.StatusBadRequest,
	)
	if customError != nil {
		c.JSON(customError.HttpStatusCode, customError)
		return
	}

	user, customError := a.service.GetUser(uint(id))
	if customError != nil {
		c.JSON(customError.HttpStatusCode, customError)
		return
	}

	responseUtils.New(user).Done(c)
}

// ShowAccount godoc
// @Summary      Create User
// @Description  Create User
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        data   body    model.UserCreate   true    "user info"
// @Success      201  {object}  response.Response
// @Router       /users/create [post]
func (a *userController) CreateUser(c *gin.Context) {
	var body model.UserCreate

	customError := errorUtils.CheckErrAndConvert(
		c.BindJSON(&body),
		http.StatusBadRequest,
	)
	if customError != nil {
		c.JSON(customError.HttpStatusCode, customError)
		return
	}

	// if body.Password != body.ConfirmPassword {
	// 	response.MakeErrorResponse(c, http.StatusBadRequest, errors.New("confirm_password doesn't match password"))
	// 	return
	// }

	user, customError := a.service.CreateUser(&body)
	if customError != nil {
		c.JSON(customError.HttpStatusCode, customError)
		return
	}

	responseUtils.New(user).SetHttpCode(http.StatusCreated).Done(c)
}

// ShowAccount godoc
// @Summary      Update User
// @Description  Update User by id
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path    int                true    "user id"
// @Param        data   body    model.UserUpdate   true    "new user info"
// @Success      202  {object}  response.Response
// @Router       /users/{id} [put]
func (a *userController) UpdateUser(c *gin.Context) {
	var updatedUser *model.UserTable
	var body model.UserUpdate

	id, err := strconv.Atoi(c.Param("id"))
	customError := errorUtils.CheckErrAndConvert(
		err,
		http.StatusBadRequest,
	)
	if customError != nil {
		c.JSON(customError.HttpStatusCode, customError)
		return
	}

	customError = errorUtils.CheckErrAndConvert(
		c.BindJSON(&body),
		http.StatusUnprocessableEntity,
	)
	if customError != nil {
		c.JSON(customError.HttpStatusCode, customError)
		return
	}

	updatedUser, customError = a.service.UpdateUser(uint(id), &body)
	if customError != nil {
		c.JSON(customError.HttpStatusCode, customError)
		return
	}

	responseUtils.New(updatedUser).SetMessage("update success").SetHttpCode(http.StatusAccepted).Done(c)
}

// ShowAccount godoc
// @Summary      Delete User
// @Description  Delete User by id
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path    int                true    "user id"
// @Success      202  {object}  response.Response
// @Router       /users/{id} [delete]
func (a *userController) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	customError := errorUtils.CheckErrAndConvert(
		err,
		http.StatusBadRequest,
	)
	if customError != nil {
		c.JSON(customError.HttpStatusCode, customError)
		return
	}

	deletedUser, customError := a.service.DeleteUser(uint(id))
	if customError != nil {
		c.JSON(customError.HttpStatusCode, customError)
		return
	}

	responseUtils.New(deletedUser).SetMessage("delete success").SetHttpCode(http.StatusAccepted).Done(c)
}
