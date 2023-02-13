package http

import (
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
	server.GET("/todos", controller.getAllTodo)
	// Get Single todo
	server.GET("/todos/:id", controller.getTodo)
	// Create todo
	server.POST("/todos/create", controller.createTodo)
	// Update todo
	server.PUT("/todos/:id", controller.updateTodo)
	// Update todo is completed
	server.PUT("/todos/:id/completed", controller.updateTodoCompleted)
	// Delete todo
	server.DELETE("/todos/:id", controller.deleteTodo)
}

// ShowAccount godoc
// @Summary             Todos
// @Description         All of todo
// @Tags                Todos
// @Accept              json
// @Produce             json
// @Security            BearerAuth
// @Param               user_id       query string false "todo id"
// @Param               title         query string false "todo title"
// @Param               is_completed  query string false "todo is completed"
// @Param               created_at    query string false "todo created at"
// @Param               updated_at    query string false "todo updated at"
// @Param               order_by      query string false "asc or desc"
// @Success             200           {object}  response.Response
// @Router              /todos        [get]
func (t *todoController) getAllTodo(c *gin.Context) {
	field := model.TodoQueries{}

	if err := c.ShouldBindQuery(&field); err != nil {
		response.MakeErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	field.UserId = c.Keys["uid"].(string)
	totalCount, allTodo, err := t.service.GetAllTodo(&field)

	if err != nil {
		response.MakeErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	field.TotalCount = *totalCount
	c.JSON(http.StatusOK, response.MakePaginationResponse(allTodo, field.Pagination))
}

// ShowAccount godoc
// @Summary      Todo
// @Description  Get todo by id
// @Tags         Todos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path    int   true    "todo id"
// @Success      200  {object}  response.Response
// @Router       /todos/{id} [get]
func (t *todoController) getTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		response.MakeErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	todo, err := t.service.GetTodo(uint(id), c.Keys["uid"].(string))

	if err != nil {
		response.MakeErrorResponse(c, http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, response.MakeCommonResponse(todo))
}

// ShowAccount godoc
// @Summary      Create Todo
// @Description  Create Todo
// @Tags         Todos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        data   body    model.TodoCreate   true    "todo info"
// @Success      201  {object}  response.Response
// @Router       /todos/create [post]
func (t *todoController) createTodo(c *gin.Context) {
	var body model.TodoCreate
	err := c.BindJSON(&body)

	if err != nil {
		response.MakeErrorResponse(c, http.StatusUnprocessableEntity, err)
		return
	}

	body.UserId = c.Keys["uid"].(string)
	err = t.service.CreateTodo(&body)

	if err != nil {
		response.MakeErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, response.MakeCommonResponse(body, http.StatusCreated))
}

// ShowAccount godoc
// @Summary      Update Todo
// @Description  Update Todo by id
// @Tags         Todos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path    int                true    "todo id"
// @Param        data   body    model.TodoUpdate   true    "new todo info"
// @Success      202    {object}  response.Response
// @Router       /todos/{id} [put]
func (t *todoController) updateTodo(c *gin.Context) {
	var id int
	var err error

	if id, err = strconv.Atoi(c.Param("id")); err != nil {
		response.MakeErrorResponse(c, http.StatusUnprocessableEntity, err)
		return
	}

	newTodo := new(model.TodoUpdate)

	if err = c.BindJSON(&newTodo); err != nil {
		response.MakeErrorResponse(c, http.StatusUnprocessableEntity, err)
		return
	}

	newTodo.Id = uint(id)
	newTodo.UserId = c.Keys["uid"].(string)
	todo, err := t.service.UpdateTodo(newTodo)

	if err != nil {
		response.MakeErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusAccepted, response.MakeCommonResponse(todo, http.StatusAccepted))
}

// ShowAccount godoc
// @Summary      Update Todo Completed
// @Description  Update Todo Completed by id
// @Tags         Todos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      int                true    "todo id"
// @Success      202    {object}  response.Response
// @Router       /todos/{id}/completed [put]
func (t *todoController) updateTodoCompleted(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.MakeErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	err = t.service.UpdateTodoCompleted(&model.TodoUpdate{
		Id:     uint(id),
		UserId: c.Keys["uid"].(string),
	})

	if err != nil {
		response.MakeErrorResponse(c, http.StatusUnprocessableEntity, err)
		return
	}

	c.JSON(http.StatusAccepted, response.MakeCommonResponse(1))
}

// ShowAccount godoc
// @Summary      Delete Todo
// @Description  Delete Todo by id
// @Tags         Todos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      int                true    "user id"
// @Success      202    {object}  response.Response
// @Router       /todos/{id} [delete]
func (t *todoController) deleteTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		response.MakeErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	deletedTodo, err := t.service.DeleteTodo(uint(id), c.Keys["uid"].(string))

	if err != nil {
		response.MakeErrorResponse(c, http.StatusUnprocessableEntity, err)
		return
	}

	c.JSON(http.StatusAccepted, response.MakeCommonResponse(deletedTodo, http.StatusAccepted))
}
