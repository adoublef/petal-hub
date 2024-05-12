package http_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	todo3 "go.petal-hub.io/gateway/internal/todo/sql3"
	. "go.petal-hub.io/gateway/net/http"
)

func Test_handle(t *testing.T) {
	t.Run("OK", run(func(t *testing.T, c *http.Client, base string) {
		// make a http request
		// read the status code
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
