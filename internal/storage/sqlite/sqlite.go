package sqlite

import (
	"database/sql"
	"log/slog"

	"github.com/MHNahib/go-rest-api-basic/internal/config"
	"github.com/MHNahib/go-rest-api-basic/internal/types"
	_ "modernc.org/sqlite"
)

type Sqlite struct {
	Db *sql.DB
}

const (
	INITIAL_QUARY = `CREATE TABLE IF NOT EXISTS todos (id INTEGER PRIMARY KEY AUTOINCREMENT, text TEXT NOT NULL, done BOOLEAN NOT NULL DEFAULT 0)`
	CREATE_QUARY  = `INSERT INTO todos (text, done) VALUES (?, ?)`
	GET_ALL_QUARY = `SELECT id, text, done FROM todos`
	UPDATE_QUARY  = `UPDATE todos SET text = ?, done = ? WHERE id = ?`
	DELETE_QUARY  = `DELETE FROM todos WHERE id = ?`
	GET_TODO      = `SELECT id, text, done FROM todos WHERE id = ?`
)

func New(config *config.Config) (*Sqlite, error) {

	db, err := sql.Open("sqlite", config.Database.Path)
	if err != nil {
		slog.Error("cannot open database", "err", err.Error())
		return nil, err
	}

	if _, err = db.Exec(INITIAL_QUARY); err != nil {
		slog.Error("cannot create table", "err", err.Error())
		return nil, err
	}

	return &Sqlite{Db: db}, nil
}

func (s *Sqlite) GetTodos() (types.Todos, error) {
	todos := []types.Todo{}

	rows, err := s.Db.Query(GET_ALL_QUARY)
	if err != nil {
		slog.Error("cannot prepare statement", "err", err.Error())
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var todo types.Todo
		if err := rows.Scan(&todo.ID, &todo.Text, &todo.Done); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (s *Sqlite) SaveTodo(todo types.Todo) (int, error) {
	statement, err := s.Db.Prepare(CREATE_QUARY)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(todo.Text, todo.Done)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s *Sqlite) DeleteTodo(id int) error {
	currentTodo, err := s.GetTodo(id)

	if err != nil {
		return err
	}

	if currentTodo.ID == 0 {
		return nil
	}

	statement, err := s.Db.Prepare(DELETE_QUARY)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Sqlite) UpdateTodo(id int, todo types.Todo) error {
	currentTodo, err := s.GetTodo(id)

	if err != nil {
		return err
	}

	if currentTodo.ID == 0 {
		return nil
	}

	statement, err := s.Db.Prepare(UPDATE_QUARY)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(todo.Text, todo.Done, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Sqlite) GetTodo(id int) (types.Todo, error) {

	if id == 0 {
		return types.Todo{}, nil
	}

	statement, err := s.Db.Prepare(GET_TODO)
	if err != nil {
		return types.Todo{}, err
	}
	defer statement.Close()

	row := statement.QueryRow(id)

	var todo types.Todo
	if err := row.Scan(&todo.ID, &todo.Text, &todo.Done); err != nil {
		if err == sql.ErrNoRows {
			return types.Todo{}, nil
		}
		return types.Todo{}, err
	}

	return todo, nil
}
