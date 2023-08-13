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

	"github.com/codescalersinternships/todoapp-Hanya/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestParseCommandLine(t *testing.T) {
	t.Run("error when no -f is specified", func(t *testing.T) {
		os.Args = []string{"./todo", "./todo.db"}
		expectedErr := ErrNoDatabaseFound
		path, port, err := ParseCommandLine()
		assert.ErrorIs(t, err, expectedErr)
		assert.Equal(t, "", path)
		assert.Equal(t, 3000, port)
	})
}

func TestInsertTodo(t *testing.T) {
	app := &App{}
	err := app.NewApp(3000, "./todo.db")
	assert.NoError(t, err, "failed to connect to database")
	app.setRoutes()
	assert.NoError(t, err, "failed to run server")
	t.Run("insertodo returns status code 201", func(t *testing.T) {
		newTodo := models.TodoItem{Task: "Test task", Completed: true}
		jsonData, err := json.Marshal(newTodo)
		assert.NoError(t, err, "failed to marshal json data")
		req, err := http.NewRequest("POST", "/todo", bytes.NewBuffer(jsonData))
		assert.NoError(t, err, "failed to create http request")
		resp := httptest.NewRecorder()
		app.router.ServeHTTP(resp, req)
		if resp.Code != http.StatusCreated {
			t.Errorf("expected status code %d but got %d", http.StatusCreated, resp.Code)
		}
	})
}

func TestGetAllTodos(t *testing.T) {
	app := &App{}
	err := app.NewApp(3000, "./todo.db")
	assert.NoError(t, err, "failed to connect to database")
	app.setRoutes()
	assert.NoError(t, err, "failed to run server")
	t.Run("getalltodos returns status 200", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/", nil)
		assert.NoError(t, err, "failed to create http request")
		resp := httptest.NewRecorder()
		app.router.ServeHTTP(resp, req)
		if resp.Code != http.StatusOK {
			t.Errorf("expected status code %d but got %d", http.StatusOK, resp.Code)
		}
	})
}

func TestGetTodo(t *testing.T) {
	app := &App{}
	err := app.NewApp(3000, "./todo.db")
	assert.NoError(t, err, "failed to connect to database")
	app.setRoutes()
	assert.NoError(t, err, "failed to run server")
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

			var req *http.Request
			if tC.errStatus == http.StatusOK {
				req, err = http.NewRequest("GET", fmt.Sprintf("/todo/%d", 1), nil)
			} else {
				req, err = http.NewRequest("GET", fmt.Sprintf("/todo/%d", -1), nil)
			}
			assert.NoError(t, err, "failed to create http request")
			resp := httptest.NewRecorder()
			app.router.ServeHTTP(resp, req)
			if resp.Code != tC.errStatus {
				t.Errorf("expected status code %d but got %d", tC.errStatus, resp.Code)
			}
		})
	}
}

func TestUpdateTodo(t *testing.T) {
	app := &App{}
	err := app.NewApp(3000, "./todo.db")
	assert.NoError(t, err, "failed to connect to database")
	app.setRoutes()
	assert.NoError(t, err, "failed to run server")
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
			newTodo := models.TodoItem{Task: "Updated task", Completed: true}
			jsonData, err := json.Marshal(newTodo)
			if err != nil {
				t.Errorf("failed to marshal json data: %v", err)
			}
			var req *http.Request
			if tC.errStatus == http.StatusOK {
				req, err = http.NewRequest("PATCH", fmt.Sprintf("/todo/%d", 1), bytes.NewBuffer(jsonData))
			} else {
				req, err = http.NewRequest("PATCH", fmt.Sprintf("/todo/%d", -1), bytes.NewBuffer(jsonData))
			}
			assert.NoError(t, err, "failed to create http request")
			resp := httptest.NewRecorder()
			app.router.ServeHTTP(resp, req)
			if resp.Code != tC.errStatus {
				t.Errorf("expected status code %d but got %d", tC.errStatus, resp.Code)
			}
		})
	}
}

func TestDeleteTodo(t *testing.T) {
	app := &App{}
	err := app.NewApp(3000, "./todo.db")
	assert.NoError(t, err, "failed to connect to database")
	app.setRoutes()
	assert.NoError(t, err, "failed to run server")
	defer os.Remove("./todo.db")
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

			var req *http.Request
			if tC.errStatus == http.StatusOK {
				req, err = http.NewRequest("DELETE", fmt.Sprintf("/todo/%d", 1), nil)
			} else {
				req, err = http.NewRequest("DELETE", fmt.Sprintf("/todo/%d", -1), nil)
			}
			assert.NoError(t, err, "failed to create http request")
			resp := httptest.NewRecorder()
			app.router.ServeHTTP(resp, req)
			if resp.Code != tC.errStatus {
				t.Errorf("expected status code %d but got %d", tC.errStatus, resp.Code)
			}
		})
	}
}
