package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/Elaman122/Go-project/internal/app/model"
	"github.com/Elaman122/Go-project/internal/app/model/filler"
	"github.com/Elaman122/Go-project/internal/jsonlog"
	vcs "github.com/Elaman122/Go-project/internal/vss"
	_ "github.com/lib/pq"
)

// config
var (
	version = vcs.Version()
)

type config struct {
	port int
	env  string
	fill bool
	db   struct {
		dsn string
	}
}

// application
type application struct {
	config config
	models model.Models
	logger *jsonlog.Logger
	wg     sync.WaitGroup
}

func main() {
	var cfg config
	flag.BoolVar(&cfg.fill, "fill", false, "Fill database with dummy data")
	flag.IntVar(&cfg.port, "port", 8081, "API server port")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgresql://postgres:ElamaN200409@@localhost/golangpro228?sslmode=disable", "PostgreSQL DSN")

	logger := jsonlog.NewLogger(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintError(err, nil)
		return
	}

	defer func() {
		if err := db.Close(); err != nil {
			logger.PrintFatal(err, nil)
		}
	}()

	app := &application{
		config: cfg,
		models: model.NewModels(db),
		logger: logger,
	}

	if cfg.fill {
		err = filler.PopulateDatabase(app.models)
		if err != nil {
			logger.PrintFatal(err, nil)
			return
		}
	}

	err = http.ListenAndServe(fmt.Sprintf(":%d", app.config.port), app.routes())
	if err != nil {
    	logger.PrintFatal(err, nil)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config // struct.
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
