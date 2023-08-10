package models

import (
	"database/sql"
	_ "embed"
	"errors"
)

type DB struct {
	client *sql.DB
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

type TodoItem struct {
	ID        string `json:"id"`
	Task      string `json:"task"`
	Completed bool   `json:"completed"`
}

var ErrTodoNotFound = errors.New("todo was not found")

func NewDB(dbPath string) (DB, error) {
	db := DB{}
	var err error
	db.client, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return db, err
	}
	_, err = db.client.Exec(createTableQuery)
	return db, err
}

func (db *DB) InsertTodo(todoItem TodoItem) error {
	_, err := db.client.Exec(insertTodoQuery, todoItem.Task, todoItem.Completed)
	return err
}
func (db *DB) GetAllTodos() ([]TodoItem, error) {
	var todoList []TodoItem
	rows, err := db.client.Query(getAllTodosQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var todoItem TodoItem
		err := rows.Scan(&todoItem.ID, &todoItem.Task, &todoItem.Completed)
		if err != nil {
			return nil, err
		}
		todoList = append(todoList, todoItem)
	}
	return todoList, nil
}

func (db *DB) GetTodo(id int) (TodoItem, error) {
	var todoItem TodoItem
	row := db.client.QueryRow(getTodoItemQuery, id)
	err := row.Scan(&todoItem.ID, &todoItem.Task, &todoItem.Completed)
	return todoItem, err
}

func (db *DB) UpdateTodo(todoItem TodoItem, id int) error {
	result, err := db.client.Exec(updateTodoItemQuery, todoItem.Task, todoItem.Completed, id)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrTodoNotFound
	}
	return err
}
func (db *DB) DeleteTodo(id int) error {
	result, err := db.client.Exec(deleteTodoItemQuery, id)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrTodoNotFound
	}
	return err
}
func (db *DB) Close() error {
	err := db.client.Close()
	return err
}
