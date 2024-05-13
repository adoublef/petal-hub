package http_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.adoublef.dev/is"
	todo3 "go.petal-hub.io/gateway/internal/todo/sql3"
	. "go.petal-hub.io/gateway/net/http"
	"go.petal-hub.io/gateway/text"
)

func Test_handleNewTodo(t *testing.T) {
	t.Run("OK", run(func(t *testing.T, c *http.Client, base string) {
		is := is.NewRelaxed(t)

		// make a http request
		// read the status code
		body := strings.NewReader(`{"title":"Learning Go!"}`)
		req, err := http.NewRequest(http.MethodPost, base+"/todo/", body)
		is.NoErr(err)

		req.Header.Add("Content-Type", "application/json")

		res, err := c.Do(req)
		is.NoErr(err)

		is.Equal(res.StatusCode, http.StatusCreated)

		loc, err := res.Location()
		is.NoErr(err)

		is.True(strings.HasPrefix(loc.Path, "/todo/"))
	}))

	t.Run("NoHeader", run(func(t *testing.T, c *http.Client, base string) {
		is := is.NewRelaxed(t)

		body := strings.NewReader(`{"title":"Learning Go!"}`)
		req, err := http.NewRequest(http.MethodPost, base+"/todo/", body)
		is.NoErr(err)

		res, err := c.Do(req)
		is.NoErr(err)

		is.Equal(res.StatusCode, http.StatusBadRequest) // missing header
	}))

	t.Run("NoTitle", run(func(t *testing.T, c *http.Client, base string) {
		is := is.NewRelaxed(t)

		body := strings.NewReader(`{"title":"Learning Go!"}`)
		req, err := http.NewRequest(http.MethodPost, base+"/todo/", body)
		is.NoErr(err)

		res, err := c.Do(req)
		is.NoErr(err)

		is.Equal(res.StatusCode, http.StatusBadRequest) // missing header
	}))
}

func Test_handleTodoList(t *testing.T) {
	t.Run("OK", run(func(t *testing.T, c *http.Client, base string) {
		is := is.NewRelaxed(t)

		// inputs
		for _, title := range []text.Title{
			"Learning Go!",
			"Learning Erlang!",
			"Learning Rust!",
		} {
			body := strings.NewReader(fmt.Sprintf(`{"title":%q}`, title))
			req, err := http.NewRequest(http.MethodPost, base+"/todo/", body)
			is.NoErr(err)

			req.Header.Add("Content-Type", "application/json")

			res, err := c.Do(req)
			is.NoErr(err)

			is.Equal(res.StatusCode, http.StatusCreated)
		}

		req, err := http.NewRequest(http.MethodGet, base+"/todo", nil)
		is.NoErr(err)

		res, err := c.Do(req)
		is.NoErr(err)

		t.Cleanup(func() {
			res.Body.Close()
		})

		// read body
		var vv []any
		is.NoErr(json.NewDecoder(res.Body).Decode(&vv))

		is.Equal(len(vv), 3)
	}))
}

func run(do func(t *testing.T, c *http.Client, base string)) func(*testing.T) {
	return func(t *testing.T) {
		db, err := todo3.Up(context.TODO(), testFilepath(t))
		if err != nil {
			t.Fatalf("new db instance: %v", err)
		}
		t.Cleanup(func() {
			db.RWC.Close()
		})

		ts := httptest.NewServer(Handler(db))
		t.Cleanup(func() {
			ts.Close()
		})

		do(t, ts.Client(), ts.URL)
	}
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func testFilepath(t testing.TB) string {
	t.Helper()

	if os.Getenv("DEBUG") != "1" {
		return filepath.Join(t.TempDir(), "test.db")
	}
	_ = os.Remove("test.db")
	return filepath.Join("test.db")
}
