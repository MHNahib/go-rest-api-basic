package storage

import "github.com/MHNahib/rest-api/internal/types"

type Storage interface {
	GetTodos() (types.Todos, error)
	SaveTodos(todo types.Todo) (int, error)
	DeleteTodo(id int) error
	UpdateTodo(id int, todo types.Todo) error
	GetTodo(id int) (types.Todo, error)
}
