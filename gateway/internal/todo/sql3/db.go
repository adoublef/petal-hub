package sql3

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
	"go.adoublef.dev/sdk/database/sql3"
	"go.adoublef.dev/sdk/time/julian"
	"go.petal-hub.io/gateway/internal/todo"
	"go.petal-hub.io/gateway/text"
)

type DB struct {
	RWC *sql3.DB
}

func (db *DB) SetTodo(ctx context.Context, title text.Title) (uuid.UUID, error) {
	created, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, wrap(err)
	}
	now := julian.Now()
	const q = `
insert into todos (id, title, body, is_complete, updated_at)
values (?, ?, ?, ?, ?)`
	_, err = db.RWC.Exec(ctx, q, created, title, nil, false, now)
	if err != nil {
		return uuid.Nil, wrap(err)
	}
	return created, nil
}

func (db *DB) Todo(ctx context.Context, tid uuid.UUID) (todo.Todo, error) {
	const q = `
select t.* from todos t where t.id = ?`
	var found Todo
	err := db.RWC.QueryRow(ctx, q, tid).Scan(&found.ID, &found.Title, &found.Body, &found.Completed, &found.Updated)
	if err != nil {
		return todo.Todo{}, wrap(err)
	}
	return TodoTo(found), nil
}

func (db *DB) Check(ctx context.Context, tid uuid.UUID) error {
	const q = `
update todos set is_complete = 1 - is_complete where id = ?`
	rs, err := db.RWC.Exec(ctx, q, tid)
	if err != nil {
		return wrap(err)
	}
	if _, err := rowsAffected(rs); err != nil {
		return wrap(err)
	}
	return nil
}

func (db *DB) SetBody(ctx context.Context, body string, tid uuid.UUID) error {
	const q = `
update todos set body = ? where id = ?`
	rs, err := db.RWC.Exec(ctx, q, body, tid)
	if err != nil {
		return wrap(err)
	}
	if _, err := rowsAffected(rs); err != nil {
		return wrap(err)
	}
	return nil
}

// NOTE TodoSlice should handle pagination
func (db *DB) TodoSlice(ctx context.Context) ([]todo.Todo, error) {
	const q = `
select t.* from todos t`
	rs, err := db.RWC.Query(ctx, q)
	if err != nil {
		return nil, wrap(err)
	}
	defer rs.Close()
	var tt []todo.Todo
	for rs.Next() {
		var t Todo
		err := rs.Scan(&t.ID, &t.Title, &t.Body, &t.Completed, &t.Updated)
		if err != nil {
			return nil, wrap(err)
		}
		tt = append(tt, TodoTo(t))
	}
	if err := rs.Err(); err != nil {
		return nil, wrap(err)
	}
	return tt, nil
}

//go:embed all:*.up.sql
var embedFS embed.FS

func Up(ctx context.Context, filename string) (*DB, error) {
	rwc, err := sql3.Up(ctx, filename, embedFS)
	if err != nil {
		return nil, fmt.Errorf("todo/sql3: running migrations: %w", err)
	}
	return &DB{rwc}, nil
}

func rowsAffected(rs sql.Result) (n int64, err error) {
	if n, err = rs.RowsAffected(); err != nil {
		return 0, err
	} else if n == 0 {
		return 0, todo.ErrNotFound
	}
	return
}

func wrap(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return errors.Join(err, todo.ErrNotFound)
	}
	// https://github.com/mattn/go-sqlite3/issues/244
	if errors.As(err, new(sqlite3.Error)) {
		switch err.(sqlite3.Error).ExtendedCode {
		case sqlite3.ErrConstraintForeignKey:
			return errors.Join(err, todo.ErrConflict)
		}
	}
	return fmt.Errorf("todo/sql3: %w", err)
}
