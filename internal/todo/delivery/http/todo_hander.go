package http

import (
	"errors"
	"net/http"
	model "simple-backend/internal/domain/todo"
	todoService "simple-backend/internal/todo/service"
	response "simple-backend/internal/utils/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type todoController struct {
	service model.TodoServiceInterface
}

func TodoHandler(server *gin.RouterGroup, DB *gorm.DB) {
	controller := &todoController{
		service: todoService.Init(DB),
	}

	if !DB.Migrator().HasTable(&model.TodoTable{}) {
		DB.AutoMigrate(&model.TodoTable{})
	}

	// Get All todo
	server.GET("/todo", controller.getAllTodo)
	// Get Single todo
	server.GET("/todo/:id", controller.getTodo)
	// Create todo
	server.POST("/todo/create", controller.createTodo)
	// Update todo
	server.PUT("/todo/:id", controller.updateTodo)
	// Update todo is completed
	server.PUT("/todo/:id/completed", controller.updateTodoCompleted)
	// Delete todo
	server.DELETE("/todo/:id", controller.deleteTodo)
}

func (t *todoController) getAllTodo(c *gin.Context) {
	field := &model.TodoQueries{}

	if err := c.ShouldBindQuery(field); err != nil {
		response.MakeErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	field.UserId = c.Keys["uid"].(string)
	totalCount, allTodo, err := t.service.GetAllTodo(field)

	if err != nil {
		response.MakeErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	field.TotalCount = *totalCount
	c.JSON(http.StatusOK, response.MakePaginationResponse(allTodo, field.Pagination))
}

func (t *todoController) getTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		response.MakeErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	todo := &model.TodoTable{}
	todo.Id = uint(id)
	todo.UserId = c.Keys["uid"].(string)
	todo, err = t.service.GetTodo(todo)

	if err != nil {
		response.MakeErrorResponse(c, http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, response.MakeCommonResponse(todo))
}

func (t *todoController) createTodo(c *gin.Context) {
	var body *model.TodoTable
	err := c.BindJSON(&body)

	if err != nil {
		response.MakeErrorResponse(c, http.StatusUnprocessableEntity, err)
		return
	}

	body.UserId = c.Keys["uid"].(string)
	err = t.service.CreateTodo(body)

	if err != nil {
		response.MakeErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, response.MakeCommonResponse(body, http.StatusCreated))
}

func (t *todoController) updateTodo(c *gin.Context) {
	var id int
	var err error

	if id, err = strconv.Atoi(c.Param("id")); err != nil {
		response.MakeErrorResponse(c, http.StatusUnprocessableEntity, err)
		return
	}

	newTodo := &model.TodoTable{}

	if err = c.BindJSON(newTodo); err != nil {
		response.MakeErrorResponse(c, http.StatusUnprocessableEntity, err)
		return
	}

	newTodo, err = t.service.UpdateTodo(id, newTodo)

	if err != nil {
		response.MakeErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusAccepted, response.MakeCommonResponse(newTodo, http.StatusAccepted))
}

func (t *todoController) updateTodoCompleted(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.MakeErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	updatedTodo := &model.TodoTable{}
	updatedTodo.Id = uint(id)
	updatedTodo.UserId = c.Keys["uid"].(string)
	updatedTodo, err = t.service.GetTodo(updatedTodo)
	if err != nil {
		response.MakeErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if updatedTodo.IsCompleted == 1 {
		response.MakeErrorResponse(c, http.StatusBadRequest, errors.New("todo is completed"))
		return
	}

	isSuccess, err := t.service.UpdateTodoCompleted(id)

	if err != nil {
		response.MakeErrorResponse(c, http.StatusUnprocessableEntity, err)
		return
	}

	c.JSON(http.StatusAccepted, response.MakeCommonResponse(isSuccess))
}

func (t *todoController) deleteTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		response.MakeErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	deletedTodo, err := t.service.DeleteTodo(id)

	if err != nil {
		response.MakeErrorResponse(c, http.StatusUnprocessableEntity, err)
		return
	}

	c.JSON(http.StatusAccepted, response.MakeCommonResponse(deletedTodo, http.StatusAccepted))
}
