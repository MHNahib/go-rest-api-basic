package todos

import (
	"encoding/json"
	"net/http"

	"github.com/MHNahib/rest-api/internal/storage"
	"github.com/MHNahib/rest-api/internal/types"
	"github.com/MHNahib/rest-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

type TodoDelete struct {
	ID int `json:"id"`
}

func GetTodosHandler(database storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		todos, err := database.GetTodos()
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GenericError(err), "")
			return
		}

		response.WriteJson(w, http.StatusOK, todos, "successfully retrieved todos")
	}
}

func CreateTodoHandler(database storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var todo types.Todo

		if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GenericError(err), "")
			return
		}

		if err := validator.New().Struct(todo); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.ValidationErrors(err.(validator.ValidationErrors)), "invalid request body")
			return
		}

		id, err := database.SaveTodo(todo)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GenericError(err), "")
			return
		}

		todo.ID = id
		response.WriteJson(w, http.StatusCreated, todo, "successfully created todo")
	}
}

func UpdateTodoHandler(database storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var todo types.Todo

		if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GenericError(err), "")
			return
		}

		if err := validator.New().Struct(todo); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.ValidationErrors(err.(validator.ValidationErrors)), "invalid request body")
			return
		}

		if todo.ID == 0 {
			response.WriteJson(w, http.StatusBadRequest, nil, "id is required")
			return
		}

		if err := database.UpdateTodo(todo.ID, todo); err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GenericError(err), "")
			return
		}

		response.WriteJson(w, http.StatusOK, todo, "successfully updated todo")
	}
}

func DeleteTodoHandler(database storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var todo TodoDelete

		if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GenericError(err), "")
			return
		}

		if todo.ID == 0 {
			response.WriteJson(w, http.StatusBadRequest, nil, "id is required")
			return
		}

		if err := database.DeleteTodo(todo.ID); err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GenericError(err), "")
			return
		}

		response.WriteJson(w, http.StatusOK, nil, "successfully deleted todo")
	}
}
