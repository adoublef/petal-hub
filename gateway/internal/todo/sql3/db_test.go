package sql3_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"go.adoublef.dev/is"
	"go.petal-hub.io/gateway/internal/todo"
	. "go.petal-hub.io/gateway/internal/todo/sql3"
	"go.petal-hub.io/gateway/text"
)

func Test_DB_SetBody(t *testing.T) {
	t.Run("OK", run(func(t *testing.T, db *DB) {
		is := is.NewRelaxed(t)

		// inputs
		var (
			title text.Title = "Learning Go!"
		)

		tid, err := db.SetTodo(context.TODO(), title)
		is.NoErr(err)

		err = db.SetBody(context.TODO(), "Hello, World!", tid)
		is.NoErr(err)

		found, err := db.Todo(context.TODO(), tid)
		is.NoErr(err)

		is.Equal(found.Body, "Hello, World!")
	}))

	t.Run("Err", run(func(t *testing.T, db *DB) {
		is := is.NewRelaxed(t)

		// inputs
		var (
			title text.Title = "Learning Go!"
		)

		_, err := db.SetTodo(context.TODO(), title)
		is.NoErr(err) // added an new todo item

		err = db.SetBody(context.TODO(), "Hello, World!", uuid.New())
		is.Err(err, todo.ErrNotFound)
	}))
}

func Test_DB_Check(t *testing.T) {
	t.Run("OK", run(func(t *testing.T, db *DB) {
		is := is.NewRelaxed(t)

		// inputs
		var (
			title text.Title = "Learning Go!"
		)

		tid, err := db.SetTodo(context.TODO(), title)
		is.NoErr(err)

		err = db.Check(context.TODO(), tid)
		is.NoErr(err) // Todo.Completed = true

		found, err := db.Todo(context.TODO(), tid)
		is.NoErr(err)

		is.Equal(found.Completed, true)

		err = db.Check(context.TODO(), tid)
		is.NoErr(err) // Todo.Completed = false

		found, err = db.Todo(context.TODO(), tid)
		is.NoErr(err)

		is.Equal(found.Completed, false)
	}))
}

func Test_DB_Todo(t *testing.T) {
	t.Run("OK", run(func(t *testing.T, db *DB) {
		is := is.NewRelaxed(t)

		// inputs
		var (
			title text.Title = "Learning Go!"
		)

		tid, err := db.SetTodo(context.TODO(), title)
		is.NoErr(err) // added an new todo item

		_, err = db.Todo(context.TODO(), tid)
		is.NoErr(err)
	}))

	t.Run("Err", run(func(t *testing.T, db *DB) {
		is := is.NewRelaxed(t)

		// inputs
		var (
			title text.Title = "Learning Go!"
		)

		_, err := db.SetTodo(context.TODO(), title)
		is.NoErr(err) // added an new todo item

		_, err = db.Todo(context.TODO(), uuid.New())
		is.Err(err, todo.ErrNotFound)
	}))
}

func Test_DB_TodoSlice(t *testing.T) {
	t.Run("OK", run(func(t *testing.T, db *DB) {
		is := is.NewRelaxed(t)

		// inputs
		for _, title := range []text.Title{
			"Learning Go!",
			"Learning Erlang!",
			"Learning Rust!",
		} {
			_, err := db.SetTodo(context.TODO(), title)
			is.NoErr(err)
		}

		tt, err := db.TodoSlice(context.TODO())
		is.NoErr(err)

		is.Equal(len(tt), 3)
	}))
}

func Test_DB_SetTodo(t *testing.T) {
	t.Run("OK", run(func(t *testing.T, db *DB) {
		is := is.NewRelaxed(t)

		// inputs
		var (
			title text.Title = "Learning Go!"
		)

		_, err := db.SetTodo(context.TODO(), title)
		is.NoErr(err)
	}))
}

func run(do func(t *testing.T, db *DB)) func(*testing.T) {
	return func(t *testing.T) {
		db, err := Up(context.TODO(), newFilepath(t))
		if err != nil {
			t.Fatalf("new db instance: %v", err)
		}
		t.Cleanup(func() {
			db.RWC.Close()
		})

		do(t, db)
	}
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func Test_Up(t *testing.T) {
	is := is.NewRelaxed(t)

	_, err := Up(context.TODO(), newFilepath(t))
	is.NoErr(err)
}

func newFilepath(t testing.TB) string {
	t.Helper()

	if os.Getenv("DEBUG") != "1" {
		return filepath.Join(t.TempDir(), "test.db")
	}
	_ = os.Remove("test.db")
	return filepath.Join("test.db")
}
