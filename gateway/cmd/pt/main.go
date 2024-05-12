package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.adoublef.dev/sdk/database/sql3"
	"go.adoublef.dev/sdk/errgroup"
	todo3 "go.petal-hub.io/gateway/internal/todo/sql3"
	handler "go.petal-hub.io/gateway/net/http"
)

func main() {
	if err := run(context.Background(), os.Args[1:], os.Getenv); err != nil {
		fmt.Fprintf(os.Stderr, "petal-hub: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, getenv func(string) string) (err error) {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	switch args[0] {
	case "serve":
		return serve(ctx, args[1:], getenv)
	case "migrate":
		return migrate(ctx, args[1:], getenv)
	default:
		return fmt.Errorf("unknown command: %q", args[0])
	}
}

func serve(ctx context.Context, args []string, _ func(string) string) (err error) {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	port := fs.Int("port", 8080, "")
	dsn := fs.String("dsn", "./todos.db", "")
	if err = fs.Parse(args); err != nil {
		return err
	}
	rwc, err := sql3.Open(*dsn)
	if err != nil {
		return err
	}
	defer rwc.Close()

	s := &http.Server{
		Addr:        fmt.Sprintf(":%d", *port),
		Handler:     handler.Handler(nil),
		BaseContext: func(l net.Listener) context.Context { return ctx },

		ReadHeaderTimeout: time.Second * 5,
		ReadTimeout:       time.Second*(10+5) + time.Millisecond*200, // read_timeout = handler_timeout + read_header_timeout + wiggle_room
		WriteTimeout:      time.Second*10 + time.Millisecond*200,     // hander_timeout + wiggle_room
		IdleTimeout:       time.Minute * 2,
	}

	return errgroup.New(ctx,
		func(ctx context.Context) (err error) {
			return s.ListenAndServe()
		},
		func(ctx context.Context) (err error) {
			<-ctx.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			// NOTE should also handle closing db connections here
			if err = s.Shutdown(ctx); err != nil {
				return err
			}
			return nil
		},
		func(ctx context.Context) (err error) {
			return http.ListenAndServe(":2112", promhttp.Handler())
		},
	).Wait()
}

func migrate(ctx context.Context, args []string, _ func(string) string) error {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	dsn := fs.String("dsn", "./todos.db", "")
	if err := fs.Parse(args); err != nil {
		return err
	}
	_, err := todo3.Up(ctx, *dsn)
	return err
}
