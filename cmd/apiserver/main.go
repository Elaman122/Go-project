package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	"github.com/Elaman122/Go-project/internal/app/model"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

//config
type config struct {
	port string
	env  string
	db   struct {
		dsn string
	}
}
// application 
type application struct {
	config config
	models model.Models
}



func main() {
	var cfg config
	flag.StringVar(&cfg.port, "port", ":8080", "API server port")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgresql://postgres:ElamaN200409@@localhost/golangpro228?sslmode=disable", "PostgreSQL DSN")

	
	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := &application{
		config: cfg,
		models: model.NewModels(db),
	}

	app.run()
}


func (app *application) run() {
    r := mux.NewRouter()

    v1 := r.PathPrefix("/api/v1").Subrouter()

    // Обработчики для создания, получения, обновления и удаления элементов меню
    v1.HandleFunc("/menus", app.createCurrencyHandler).Methods("POST")
    v1.HandleFunc("/menus/{menuId:[0-9]+}", app.getCurrencyHandler).Methods("GET")
    v1.HandleFunc("/menus/{menuId:[0-9]+}", app.updateCurrencyHandler).Methods("PUT")
    v1.HandleFunc("/menus/{menuId:[0-9]+}", app.deleteCurrencyHandler).Methods("DELETE")

    // Обработчик для получения списка меню с поддержкой пагинации, сортировки и фильтрации
    v1.HandleFunc("/menufor", app.getAllCurrenciesHandler).Methods("GET")


    log.Printf("Starting server on %s\n", app.config.port)
    err := http.ListenAndServe(app.config.port, r)
    log.Fatal(err)
}

func openDB(cfg config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config // struct.
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
