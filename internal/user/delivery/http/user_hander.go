package http

import (
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
	users, err := a.service.GetUsers()

	if err != nil {
		response.MakeErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, response.MakeCommonResponse(users))
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
	if err := c.BindJSON(&body); err != nil {
		response.MakeErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	// if body.Password != body.ConfirmPassword {
	// 	response.MakeErrorResponse(c, http.StatusBadRequest, errors.New("confirm_password doesn't match password"))
	// 	return
	// }

	user, err := a.service.CreateUser(&body)
	if err != nil {
		response.MakeErrorResponse(c, http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusCreated, response.MakeCommonResponse(user, http.StatusCreated))
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
	var id int
	var body model.UserUpdate
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
