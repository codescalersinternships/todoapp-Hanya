package todo

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"strconv"

	"github.com/codescalersinternships/todoapp-Hanya/internal/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type App struct {
	db     models.DB
	router *gin.Engine
}

var (
	// ErrNoDatabaseFound indicates abscence of -f flag
	ErrNoDatabaseFound = errors.New("no database path specified")
)

func ParseCommandLine() (string, int, error) {
	var dbPath string
	var port int
	flag.StringVar(&dbPath, "f", "", "path to the database file")
	flag.IntVar(&port, "p", 3000, "localhost port for the server")
	flag.Parse()
	if dbPath == "" {
		return "", port, ErrNoDatabaseFound
	}
	return dbPath, port, nil
}

func (app *App) NewApp(port int, dbPath string) error {
	db, err := models.NewDB(dbPath)
	app.db = db
	return err
}

// @Summary      inserts new task
// @Description  accepts task's title and completed bool in json format and returns the created task in json
// @Produce      json
// @Tags         create
// @Param        item  body      models.TodoItem  true  "task body"
// @Success      201   {object}  models.TodoItem
// @Router       /todo/ [post]
// @Failure      400   "error binding json data"
// @Failure      500   "error inserting new task"
func (app *App) insertTodo(c *gin.Context) (interface{}, Response) {
	var todoItem models.TodoItem
	if err := c.ShouldBindJSON(&todoItem); err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("error binding json data"))
	}
	err := app.db.InsertTodo(todoItem)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errors.New("error inserting new task"))
	}
	return ResponseMsg{
		Message: "task inserted successfully",
		Data:    todoItem,
	}, Created()
}

// @Summary      gets all tasks
// @Description  retrieves a list of all tasks in database in json format
// @Produce      json
// @Tags         get
// @Success      200  {object}  []models.TodoItem
// @Router       / [get]
// @Failure      400   "error retrieving all tasks"
func (app *App) getAllTodos(c *gin.Context) (interface{}, Response) {
	todoList, err := app.db.GetAllTodos()
	if err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("error retrieving all tasks"))
	}
	return ResponseMsg{
		Message: "all tasks retrieved successfully",
		Data:    todoList,
	}, Ok()
}

// @Summary      gets a task by id
// @Description  retrieves a list of all tasks in database in json format
// @Produce      json
// @Tags         get
// @Param        id  path      string  true  "required task's id"
// @Success      200  {object}  models.TodoItem
// @Router       /todo/:id [get]
// @Failure      400   "error retrieving task"
func (app *App) getTodo(c *gin.Context) (interface{}, Response) {
	id, _ := strconv.Atoi(c.Param("id"))
	todoItem, err := app.db.GetTodo(id)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, NotFound(errors.New("error retrieving task"))
	}
	return ResponseMsg{
		Message: "task retrieved successfully",
		Data:    todoItem,
	}, Ok()
}

// @Summary      updates a task
// @Description  accepts task's title and completed bool in json format and returns the created task in json
// @Tags	 	 update
// @Produce      json
// @Param        item  body      models.TodoItem true "required task's id"
// @Success      200  {object}  string
// @Router       /todo/:id [patch]
// @Failure      400  "error binding json data"
// @Failure      500  "error updating task"
func (app *App) updateTodo(c *gin.Context) (interface{}, Response) {
	var todoItem models.TodoItem
	if err := c.ShouldBindJSON(&todoItem); err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("error binding json data"))
	}
	id, _ := strconv.Atoi(c.Param("id"))
	err := app.db.UpdateTodo(todoItem, id)
	if err != nil {
		log.Error().Err(err).Send()
		if err == models.ErrTodoNotFound {
			return nil, NotFound(err)
		}
		return nil, InternalServerError(errors.New("error updating task"))
	}
	return ResponseMsg{
		Message: "task updated successfully",
		Data:    todoItem,
	}, Ok()
}

// @Summary      deletes a task
// @Description  accepts task's id and deletes it
// @Tags	 	 delete
// @Produce      json
// @Param        item  body      models.TodoItem true "required task's id"
// @Success      200  {object}  string
// @Router       /todo/:id [delete]
// @Failure      500  "error deleting task"
func (app *App) deleteTodo(c *gin.Context) (interface{}, Response) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := app.db.DeleteTodo(id)
	if err != nil {
		if err == models.ErrTodoNotFound {
			log.Error().Err(err).Send()
			return nil, NotFound(err)

		}
		return nil, InternalServerError(errors.New("error deleting task"))
	}
	return ResponseMsg{
		Message: "task deleted successfully",
	}, Ok()
}
func (app *App) Run(port int) error {
	app.router = gin.Default()
	app.setRoutes()
	return app.router.Run(fmt.Sprintf(":%d", port))
}
func (app *App) setRoutes() {
	app.router = gin.Default()
	app.router.Use(cors.Default())
	app.router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	app.router.POST("/todo", WrapFunc(app.insertTodo))
	app.router.GET("/", WrapFunc(app.getAllTodos))
	app.router.GET("/todo/:id", WrapFunc(app.getTodo))
	app.router.PATCH("/todo/:id", WrapFunc(app.updateTodo))
	app.router.DELETE("/todo/:id", WrapFunc(app.deleteTodo))

}
func (app *App) Close() error {
	err := app.db.Close()
	return err
}
