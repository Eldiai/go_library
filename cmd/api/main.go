package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Eldiai/go_library/config"
	"github.com/Eldiai/go_library/internal/jsonlog"
	"github.com/Eldiai/go_library/internal/mailer"

	_ "github.com/jackc/pgx/v5/stdlib" // for compatibility with database/sql
)

const version = "1.0.0"

type application struct {
	config *config.Config
	logger *jsonlog.Logger
	mailer mailer.Mailer
}

func main() {
	cfg := config.GetConfig()

	app := &application{
		config: cfg,
		logger: jsonlog.NewLogger(os.Stdout, jsonlog.LevelInfo),
		mailer: mailer.New(cfg.Smtp.Host, cfg.Smtp.Port, cfg.Smtp.Username, cfg.Smtp.Password, cfg.Smtp.Sender),
	}

	srv := &http.Server{
		Addr:         cfg.Port,
		Handler:      app.routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			app.logger.PrintFatal(err, nil)
		}
	}()

	// gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		app.logger.PrintFatal(err, nil)
	}
}

func openDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.Db.Dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.Db.MaxOpenConns)

	db.SetMaxIdleConns(cfg.Db.MaxIdleConns)

	duration, err := time.ParseDuration(cfg.Db.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	// Set the maximum idle timeout.
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
