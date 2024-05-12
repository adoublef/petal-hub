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

func JournalOf(t todo.Todo) Todo {
	var ptr *string
	if t.Body != "" {
		ptr = &t.Body
	}
	return Todo{t.ID, t.Title, ptr, t.Completed, julian.FromTime(t.Updated)}
}
