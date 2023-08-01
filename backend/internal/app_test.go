package todo

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestParseCommandLine(t *testing.T) {
	t.Run("error when no -f is specified", func(t *testing.T) {
		os.Args = []string{"./todo", "../todo.db"}
		expectedErr := ErrNoDatabaseFound
		path, port, err := ParseCommandLine()
		assert.ErrorIs(t, err, expectedErr)
		assert.Equal(t, "", path)
		assert.Equal(t,3000,port)
	})
}

func TestNewApp(t *testing.T) {
	t.Run("newdb returns no error", func(t *testing.T) {
		app := &App{DbPath: "../todo.db"}
		err := app.NewApp(3000)
		assert.NoError(t, err)
	})
}

func TestInsertTodo(t *testing.T) {
	t.Run("insertodo returns status code 201", func(t *testing.T) {
		router := gin.Default()
		app := &App{DbPath: "../todo.db"}
		err := app.NewApp(3000)
		assert.NoError(t, err, "failed to connect to database")
		router.POST("/todos", app.insertTodo)
		newTodo := TodoItem{Task: "Test task", Completed: true}
		jsonData, err := json.Marshal(newTodo)
		assert.NoError(t, err, "failed to marshal json data")
		req, err := http.NewRequest("POST", "/todos", bytes.NewBuffer(jsonData))
		assert.NoError(t, err, "failed to create http request")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		if resp.Code != http.StatusCreated {
			t.Errorf("expected status code %d but got %d", http.StatusCreated, resp.Code)
		}
	})
}

func TestGetAllTodos(t *testing.T) {
	t.Run("getalltodos returns status 200", func(t *testing.T) {
		router := gin.Default()
		app := &App{DbPath: "../todo.db"}
		err := app.NewApp(3000)
		assert.NoError(t, err, "failed to connect to database")
		router.GET("/", app.getAllTodos)
		req, err := http.NewRequest("GET", "/", nil)
		assert.NoError(t, err, "failed to create http request")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		if resp.Code != http.StatusOK {
			t.Errorf("expected status code %d but got %d", http.StatusOK, resp.Code)
		}
	})
}

//go:embed sql/01getFirstTodoItemId.sql
var getTodoItemTestQuery string

func TestGetTodo(t *testing.T) {
	testCases := []struct {
		name      string
		errStatus int
	}{
		{
			name:      "gettodo returns status 200",
			errStatus: http.StatusOK,
		},
		{
			name:      "gettodo returns status 404",
			errStatus: http.StatusNotFound,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			router := gin.Default()
			app := &App{DbPath: "../todo.db"}
			err := app.NewApp(3000)
			assert.NoError(t, err, "failed to connect to database")
			var id int
			err = app.db.QueryRow(getTodoItemTestQuery).Scan(&id)
			assert.NoError(t, err, "failed to get id of first row")
			router.GET("/todo/:id", app.getTodo)
			var req *http.Request
			if tC.errStatus == http.StatusOK {
				req, err = http.NewRequest("GET", fmt.Sprintf("/todo/%d", id), nil)
			} else {
				req, err = http.NewRequest("GET", fmt.Sprintf("/todo/%d", id-1), nil)
			}
			assert.NoError(t, err, "failed to create http request")
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			if resp.Code != tC.errStatus {
				t.Errorf("expected status code %d but got %d", tC.errStatus, resp.Code)
			}
		})
	}
}

func TestUpdateTodo(t *testing.T) {
	testCases := []struct {
		name      string
		errStatus int
	}{
		{
			name:      "updatetodo returns status 200",
			errStatus: http.StatusOK,
		},
		{
			name:      "updatetodo returns status 404",
			errStatus: http.StatusNotFound,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			router := gin.Default()
			app := &App{DbPath: "../todo.db"}
			err := app.NewApp(3000)
			assert.NoError(t, err, "failed to connect to database")
			var id int
			err = app.db.QueryRow(getTodoItemTestQuery).Scan(&id)
			fmt.Println("IDDDDD",id)
			assert.NoError(t, err, "failed to get id of first row")
			router.PATCH("/todo/:id", app.updateTodo)
			newTodo := TodoItem{Task: "Updated task", Completed: true}
			jsonData, err := json.Marshal(newTodo)
			if err != nil {
				t.Errorf("failed to marshal json data: %v", err)
			}
			var req *http.Request
			if tC.errStatus == http.StatusOK {
				req, err = http.NewRequest("PATCH", fmt.Sprintf("/todo/%d", id), bytes.NewBuffer(jsonData))
			} else {
				req, err = http.NewRequest("PATCH", fmt.Sprintf("/todo/%d", id-1), bytes.NewBuffer(jsonData))
			}
			assert.NoError(t, err, "failed to create http request")
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			if resp.Code != tC.errStatus {
				t.Errorf("expected status code %d but got %d", tC.errStatus, resp.Code)
			}
		})
	}
}

func TestDeleteTodo(t *testing.T) {
	testCases := []struct {
		name      string
		errStatus int
	}{
		{
			name:      "updatetodo returns status 200",
			errStatus: http.StatusOK,
		},
		{
			name:      "updatetodo returns status 404",
			errStatus: http.StatusNotFound,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			router := gin.Default()
			app := &App{DbPath: "../todo.db"}
			err := app.NewApp(3000)
			assert.NoError(t, err, "failed to connect to database")
			var id int
			err = app.db.QueryRow(getTodoItemTestQuery).Scan(&id)
			assert.NoError(t, err, "failed to get id of first row")
			router.DELETE("/todo/:id", app.deleteTodo)
			var req *http.Request
			if tC.errStatus == http.StatusOK {
				req, err = http.NewRequest("DELETE", fmt.Sprintf("/todo/%d", id), nil)
			} else {
				req, err = http.NewRequest("DELETE", fmt.Sprintf("/todo/%d", id-1), nil)
			}
			assert.NoError(t, err, "failed to create http request")
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			if resp.Code != tC.errStatus {
				t.Errorf("expected status code %d but got %d", tC.errStatus, resp.Code)
			}
		})
	}
}

func TestClose(t *testing.T) {
	t.Run("close ends database connection", func(t *testing.T) {
		app := &App{DbPath: "../todo.db"}
		err := app.NewApp(3000)
		assert.NoError(t, err)
		err = app.Close()
		assert.NoError(t, err)
	})
}
