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
	ID    uuid.UUID
	Title text.Title
	// TODO body has a constraint to not exceed 280 characters
	Body      string
	Completed bool
	Updated   time.Time
}
