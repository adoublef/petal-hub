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

		err = db.SetBody(context.TODO(), "Hello, World!", tid, 0)
		is.NoErr(err)

		found, _, err := db.Todo(context.TODO(), tid)
		is.NoErr(err)

		is.Equal(found.Body, "Hello, World!")
	}))

	t.Run("NotFound", run(func(t *testing.T, db *DB) {
		is := is.NewRelaxed(t)

		// inputs
		var (
			title text.Title = "Learning Go!"
		)

		_, err := db.SetTodo(context.TODO(), title)
		is.NoErr(err) // added an new todo item

		err = db.SetBody(context.TODO(), "Hello, World!", uuid.New(), 0)
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

		err = db.Check(context.TODO(), tid, 0) // v=0
		is.NoErr(err)                          // Todo.Completed = true

		found, v, err := db.Todo(context.TODO(), tid)
		is.NoErr(err)

		is.Equal(found.Completed, true)

		err = db.Check(context.TODO(), tid, v)
		is.NoErr(err) // Todo.Completed = false

		found, _, err = db.Todo(context.TODO(), tid)
		is.NoErr(err)

		is.Equal(found.Completed, false)
	}))

	t.Run("NotFound", run(func(t *testing.T, db *DB) {
		is := is.NewRelaxed(t)

		// inputs
		var (
			title text.Title = "Learning Go!"
		)

		tid, err := db.SetTodo(context.TODO(), title)
		is.NoErr(err)

		err = db.Check(context.TODO(), tid, 1) // v=0
		is.Err(err, todo.ErrNotFound)
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

		_, _, err = db.Todo(context.TODO(), tid)
		is.NoErr(err)
	}))

	t.Run("NotFound", run(func(t *testing.T, db *DB) {
		is := is.NewRelaxed(t)

		// inputs
		var (
			title text.Title = "Learning Go!"
		)

		_, err := db.SetTodo(context.TODO(), title)
		is.NoErr(err) // added an new todo item

		_, _, err = db.Todo(context.TODO(), uuid.New())
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

		tt, err := db.Slice(context.TODO())
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
		db, err := Up(context.TODO(), testFilepath(t))
		if err != nil {
			t.Fatalf("new db instance: %v", err)
		}
		t.Cleanup(func() {
			db.RWC.Close()
		})

		do(t, db)
	}
}

func Test_Up(t *testing.T) {
	is := is.NewRelaxed(t)

	_, err := Up(context.TODO(), testFilepath(t))
	is.NoErr(err)
}

func testFilepath(t testing.TB) string {
	t.Helper()

	if os.Getenv("DEBUG") != "1" {
		return filepath.Join(t.TempDir(), "test.db")
	}
	_ = os.Remove("test.db")
	return filepath.Join("test.db")
}
