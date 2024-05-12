package http

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	lop "github.com/samber/lo/parallel"
	"go.petal-hub.io/gateway/internal/todo"
	todo3 "go.petal-hub.io/gateway/internal/todo/sql3"
	"go.petal-hub.io/gateway/text"
)

func Handler(todoDB *todo3.DB) http.Handler {
	mux := http.NewServeMux()
	applyRoutes(mux, todoDB)
	return mux
}

func applyRoutes(mux *http.ServeMux, todoDB *todo3.DB) {
	mux.Handle("GET /{$}", http.RedirectHandler("/todo/", http.StatusFound))
	mux.HandleFunc("GET /todo/{$}", handleTodoList(todoDB))
	mux.HandleFunc("POST /todo/{$}", handleNewTodo(todoDB))
	mux.HandleFunc("GET /todo/{todoID}", handleTodo(todoDB))
	mux.HandleFunc("PUT /todo/{todoID}/body", handleUpdateBody(todoDB))
	mux.HandleFunc("PUT /todo/{todoID}/complete", handleCheckTodo(todoDB))
}

func handleNewTodo(todoDB *todo3.DB) http.HandlerFunc {
	type request struct {
		Title *text.Title `json:"title"`
	}
	var parse = func(w http.ResponseWriter, r *http.Request) (text.Title, error) {
		rq, err := decode[request](w, r)
		if err != nil {
			return "", err
		}
		if rq.Title == nil {
			return "", fmt.Errorf("http: title is required")
		}
		return *rq.Title, nil
	}
	return func(w http.ResponseWriter, r *http.Request) {
		title, err := parse(w, r)
		if err != nil {
			http.Error(w, "Invalid Title", http.StatusBadRequest)
			return
		}
		created, err := todoDB.SetTodo(r.Context(), title)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		w.Header().Add("Location", fmt.Sprintf("/todo/%s", created))
		w.WriteHeader(http.StatusCreated)
	}
}

func handleTodo(todoDB *todo3.DB) http.HandlerFunc {
	var parse = func(_ http.ResponseWriter, r *http.Request) (uuid.UUID, error) {
		tid, err := uuid.Parse(r.PathValue("todoID"))
		if err != nil {
			return uuid.Nil, fmt.Errorf("http: parsing path id: %w", err)
		}
		// should be able to tell the version
		return tid, nil
	}
	type response struct {
		Todo    todo.Todo `json:"todo"`
		Version uint      `json:"version"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		tid, err := parse(w, r)
		if err != nil {
			// return a 404 is a personal preference
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		found, v, err := todoDB.Todo(r.Context(), tid)
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		encode(w, r, http.StatusOK, &response{found, v})
	}
}

func handleCheckTodo(todoDB *todo3.DB) http.HandlerFunc {
	type request struct {
		Version uint `json:"version"`
	}
	var parse = func(w http.ResponseWriter, r *http.Request) (uuid.UUID, uint, error) {
		tid, err := uuid.Parse(r.PathValue("todoID"))
		if err != nil {
			return uuid.Nil, 0, fmt.Errorf("http: parsing path id: %w", err)
		}
		rq, err := decode[request](w, r)
		if err != nil {
			return uuid.Nil, 0, err
		}
		return tid, rq.Version, nil
	}
	return func(w http.ResponseWriter, r *http.Request) {
		tid, v, err := parse(w, r)
		if err != nil {
			// return a 404 is a personal preference
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		err = todoDB.Check(r.Context(), tid, v)
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func handleUpdateBody(todoDB *todo3.DB) http.HandlerFunc {
	type request struct {
		Body    string `json:"body"`
		Version uint   `json:"version"`
	}
	var parse = func(w http.ResponseWriter, r *http.Request) (string, uuid.UUID, uint, error) {
		tid, err := uuid.Parse(r.PathValue("todoID"))
		if err != nil {
			return "", uuid.Nil, 0, fmt.Errorf("http: parsing path id: %w", err)
		}
		rq, err := decode[request](w, r)
		if err != nil {
			return "", uuid.Nil, 0, err
		}
		return rq.Body, tid, rq.Version, nil
	}
	return func(w http.ResponseWriter, r *http.Request) {
		body, tid, v, err := parse(w, r)
		if err != nil {
			// return a 404 is a personal preference
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		err = todoDB.SetBody(r.Context(), body, tid, v)
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// handle pagination
func handleTodoList(todoDB *todo3.DB) http.HandlerFunc {
	type response struct {
		Todo todo.Summary `json:"todo"`
		URL  string       `json:"urlPath"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// NOTE this would be better with pagination
		// and using an iterator rather than a slice
		tt, err := todoDB.Slice(r.Context())
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		rr := lop.Map(tt, func(t todo.Summary, i int) response {
			return response{t, fmt.Sprintf("/todo/%s", t.ID)}
		})
		encode(w, r, http.StatusOK, rr)
	}
}
