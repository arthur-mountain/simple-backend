package http

import (
	"net/http"
	model "simple-backend/internal/domain/todo"
	todoService "simple-backend/internal/todo/service"
	"simple-backend/internal/utils/databases"
	errorUtils "simple-backend/internal/utils/error"
	responseUtils "simple-backend/internal/utils/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type todoController struct {
	service model.TodoServiceInterface
}

func TodoHandler(server *gin.RouterGroup, db *databases.TMysql) {
	controller := &todoController{
		service: todoService.Init(db),
	}

	db.Execute(func(DB *gorm.DB) error {
		table := new(model.TodoTable)
		if !DB.Migrator().HasTable(table) {
			DB.AutoMigrate(table)
		}
		return nil
	}, nil)

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

	customError := errorUtils.CheckErrAndConvert(
		c.ShouldBindQuery(&field),
		http.StatusBadRequest,
	)
	if customError != nil {
		c.JSON(customError.HttpStatusCode, customError)
		return
	}

	field.UserId = c.Keys["uid"].(string)
	totalCount, allTodo, customError := t.service.GetAllTodo(&field)

	if customError != nil {
		c.JSON(customError.HttpStatusCode, customError)
		return
	}

	field.TotalCount = *totalCount
	responseUtils.New(allTodo).AppendPagination(&field.Pagination).Done(c)
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

	customError := errorUtils.CheckErrAndConvert(
		err,
		http.StatusBadRequest,
	)
	if customError != nil {
		c.JSON(customError.HttpStatusCode, customError)
		return
	}

	todo, customError := t.service.GetTodo(uint(id), c.Keys["uid"].(string))

	if customError != nil {
		c.JSON(customError.HttpStatusCode, customError)
		return
	}

	responseUtils.New(todo).Done(c)
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

	customError := errorUtils.CheckErrAndConvert(
		c.BindJSON(&body),
		http.StatusBadRequest,
	)
	if customError != nil {
		c.JSON(customError.HttpStatusCode, customError)
		return
	}

	body.UserId = c.Keys["uid"].(string)
	customError = t.service.CreateTodo(&body)

	if customError != nil {
		c.JSON(customError.HttpStatusCode, customError)
		return
	}

	responseUtils.New(body).SetHttpCode(http.StatusCreated).Done(c)
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
	id, err := strconv.Atoi(c.Param("id"))

	customError := errorUtils.CheckErrAndConvert(
		err,
		http.StatusUnprocessableEntity,
	)
	if customError != nil {
		c.JSON(customError.HttpStatusCode, customError)
		return
	}

	newTodo := new(model.TodoUpdate)

	customError = errorUtils.CheckErrAndConvert(
		c.BindJSON(&newTodo),
		http.StatusUnprocessableEntity,
	)
	if customError != nil {
		c.JSON(customError.HttpStatusCode, customError)
		return
	}

	newTodo.Id = uint(id)
	newTodo.UserId = c.Keys["uid"].(string)
	todo, customError := t.service.UpdateTodo(newTodo)

	if customError != nil {
		c.JSON(customError.HttpStatusCode, customError)
		return
	}

	responseUtils.New(todo).SetHttpCode(http.StatusAccepted).Done(c)
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

	customError := errorUtils.CheckErrAndConvert(
		err,
		http.StatusUnprocessableEntity,
	)
	if customError != nil {
		c.JSON(customError.HttpStatusCode, customError)
		return
	}

	customError = t.service.UpdateTodoCompleted(&model.TodoUpdate{
		Id:     uint(id),
		UserId: c.Keys["uid"].(string),
	})

	if customError != nil {
		c.JSON(customError.HttpStatusCode, customError)
		return
	}

	responseUtils.New(1).SetHttpCode(http.StatusAccepted).Done(c)
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

	customError := errorUtils.CheckErrAndConvert(
		err,
		http.StatusUnprocessableEntity,
	)
	if customError != nil {
		c.JSON(customError.HttpStatusCode, customError)
		return
	}

	deletedTodo, customError := t.service.DeleteTodo(uint(id), c.Keys["uid"].(string))

	if customError != nil {
		c.JSON(customError.HttpStatusCode, customError)
		return
	}

	responseUtils.New(deletedTodo).SetHttpCode(http.StatusAccepted).Done(c)
}
