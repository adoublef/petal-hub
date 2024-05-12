package sql3

import (
	"github.com/google/uuid"
	"go.adoublef.dev/sdk/time/julian"
	"go.petal-hub.io/gateway/internal/todo"
	"go.petal-hub.io/gateway/text"
)

type Todo struct {
	ID        uuid.UUID
	Title     text.Title
	Body      *string
	Completed bool
	Updated   julian.Time
}

func TodoTo(t Todo) todo.Todo {
	var body string
	if t.Body != nil {
		body = *t.Body
	}
	return todo.Todo{t.ID, t.Title, body, t.Completed, t.Updated.Time()}
}

type Summary struct {
	ID        uuid.UUID
	Title     text.Title
	Completed bool
	Updated   julian.Time
}

func SummaryTo(s Summary) todo.Summary {
	return todo.Summary{s.ID, s.Title, s.Completed, s.Updated.Time()}
}
