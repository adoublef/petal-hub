package todo

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"go.petal-hub.io/gateway/text"
)

var (
	ErrNotFound = errors.New("todo: not found")
	ErrConflict = errors.New("todo: conflict")
)

type Todo struct {
	ID    uuid.UUID  `json:"id"`
	Title text.Title `json:"title"`
	// TODO body has a constraint to not exceed 280 characters
	Body      string    `json:"body"`
	Completed bool      `json:"isComplete"`
	Updated   time.Time `json:"updatedAt"`
}

// Summary is like [Todo] excluding the body as this can be quite large
type Summary struct {
	ID        uuid.UUID  `json:"id"`
	Title     text.Title `json:"title"`
	Completed bool       `json:"isComplete"`
	Updated   time.Time  `json:"updatedAt"`
}
