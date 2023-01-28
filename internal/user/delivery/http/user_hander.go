package http

import (
	"errors"
	model "simple-backend/internal/domain/user"
	userService "simple-backend/internal/user/service"
	response "simple-backend/internal/utils/response"
	"strconv"

	"simple-backend/internal/utils/databases"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type userController struct {
	service model.UserServiceInterface
}

func UserHandler(server *gin.RouterGroup, DB *gorm.DB, REDIS *databases.MyRedis) {
	controller := &userController{
		service: userService.Init(DB, REDIS),
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

func (a *userController) GetUsers(c *gin.Context) {
	users, err := a.service.GetUsers()

	if err != nil {
		response.MakeErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, response.MakeCommonResponse(users))
}

// Get User
func (a *userController) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.MakeErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	user, err := a.service.GetUser(uint(id))
	if err != nil {
		response.MakeErrorResponse(c, http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, response.MakeCommonResponse(user))
}

// Create User
func (a *userController) CreateUser(c *gin.Context) {
	var body model.UserBody
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
func (a *userController) UpdateUser(c *gin.Context) {
	var id int
	var body model.UserBody
	var err error

	if id, err = strconv.Atoi(c.Param("id")); err != nil {
		response.MakeErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err = c.BindJSON(&body); err != nil {
		response.MakeErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	err = a.service.UpdateUser(uint(id), &body)
	if err != nil {
		response.MakeErrorResponse(c, http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusAccepted, response.MakeCommonResponse("update success", http.StatusAccepted))
}

// Delete User
func (a *userController) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.MakeErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	err = a.service.DeleteUser(uint(id))
	if err != nil {
		response.MakeErrorResponse(c, http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusAccepted, response.MakeCommonResponse("delete success", http.StatusAccepted))
}
