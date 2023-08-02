package todo

import (
	"database/sql"
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	db     *sql.DB
	DbPath string
	port   int
}

type TodoItem struct {
	ID        string `json:"id"`
	Task      string `json:"task"`
	Completed bool   `json:"completed"`
}

//go:embed sql/01createTableTodo.sql
var createTableQuery string

//go:embed sql/01insertTodoItem.sql
var insertTodoQuery string

//go:embed sql/01getAllTodoList.sql
var getAllTodosQuery string

//go:embed sql/01getTodoItem.sql
var getTodoItemQuery string

//go:embed sql/01updateTodoItem.sql
var updateTodoItemQuery string

//go:embed sql/01deleteTodoItem.sql
var deleteTodoItemQuery string

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

func (app *App) NewApp(port int) error {
	database, err := sql.Open("sqlite3", app.DbPath)
	if err != nil {
		return err
	}
	app.db = database
	app.port = port
	_, err = app.db.Exec(createTableQuery)
	return err
}

func (app *App) insertTodo(c *gin.Context) {
	var todoItem TodoItem
	if err := c.ShouldBindJSON(&todoItem); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	_, err := app.db.Exec(insertTodoQuery, todoItem.Task, todoItem.Completed)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusCreated, todoItem)
}

func (app *App) getAllTodos(c *gin.Context) {
	var todoList []TodoItem
	rows, err := app.db.Query(getAllTodosQuery)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var todoItem TodoItem
		err := rows.Scan(&todoItem.ID, &todoItem.Task, &todoItem.Completed)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		todoList = append(todoList, todoItem)
	}
	c.JSON(http.StatusOK, todoList)
}

func (app *App) getTodo(c *gin.Context) {
	var todoItem TodoItem
	row := app.db.QueryRow(getTodoItemQuery, c.Param("id"))
	err := row.Scan(&todoItem.ID, &todoItem.Task, &todoItem.Completed)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, todoItem)
}

func (app *App) updateTodo(c *gin.Context) {
	var todoItem TodoItem
	if err := c.ShouldBindJSON(&todoItem); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	result, err := app.db.Exec(updateTodoItemQuery, todoItem.Task, boolToInt(todoItem.Completed), c.Param("id"))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, todoItem)
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (app *App) deleteTodo(c *gin.Context) {
	result, err := app.db.Exec(deleteTodoItemQuery, c.Param("id"))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		c.Status(http.StatusNotFound)
		return
	}
	c.Status(http.StatusOK)
}
func (app *App) Run() error {
	router := gin.New()
	router.Use(cors.Default())
	router.POST("/todo", app.insertTodo)
	router.GET("/", app.getAllTodos)
	router.GET("/todo/:id", app.getTodo)
	router.PATCH("/todo/:id", app.updateTodo)
	router.DELETE("/todo/:id", app.deleteTodo)
	err := http.ListenAndServe(fmt.Sprintf(":%d", app.port), router)
	return err
}

func (app *App) Close() error {
	err := app.db.Close()
	return err
}
