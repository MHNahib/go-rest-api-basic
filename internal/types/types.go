package types

type Todo struct {
	ID   int    `json:"id"`
	Text string `json:"text" validate:"required"`
	Done bool   `json:"done"`
}

type Todos []Todo
