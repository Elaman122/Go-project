package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"time"

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

type Menu struct {
    ID        int       `json:"id"`
    Code      string    `json:"code"`
    Rate      float64   `json:"rate"`
    Timestamp time.Time `json:"timestamp"`
}


func main() {
	var cfg config
	flag.StringVar(&cfg.port, "port", ":8080", "API server port")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgresql://postgres:ElamaN200409@@localhost/golangpro3?sslmode=disable", "PostgreSQL DSN")

	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := &application{
		config: cfg,
		models: model.NewModels(db),
	}

	menus := []model.Menu{
		{ID: 64, Code: "USD", Rate: 1.0},
		{ID: 65, Code: "EUR", Rate: 0.84},
		{ID: 66, Code: "GBP", Rate: 1.23},
		{ID: 67, Code: "RUB", Rate: 0.011},
		{ID: 68, Code: "KZT", Rate: 0.0023},
	}

	
	// Используем err для операции Insert
	/*
	for _, m := range menus {
		if err := app.models.Menu.Insert(&m); err != nil {
			log.Printf("Ошибка при добавлении элемента в базу данных: %v", err)
			} else {
				log.Println("Элемент успешно добавлен в базу данных")
			}
	}
	*/	
	// Используем err для операции Update
	updatedMenu := &model.Menu{
		ID:        60,
		Code:      "NewCode",
		Rate:      5.0,
		Timestamp: time.Now(),
	}

    //update

	if err := app.models.Menu.Update(updatedMenu); err != nil {
		log.Printf("Ошибка при обновлении элемента в базе данных: %v", err)
	} else {
		log.Println("Элемент успешно обновлен в базу данных")
	}


    // Delete
    
    for _, menu := range menus {
		if err := app.models.Menu.Delete(menu.ID); err != nil {
			log.Printf("Ошибка при удалении элемента из базы данных: %v", err)
		} else {
			log.Printf("Элемент с ID %d успешно удален из базы данных", menu.ID)
		}
	}
	app.run()
}

func (app *application) run() {
	r := mux.NewRouter()

	v1 := r.PathPrefix("/api/v1").Subrouter()

	v1.HandleFunc("/menus", app.createCurrencyHandler).Methods("POST")
	v1.HandleFunc("/menus/{menuId:[0-9]+}", app.getCurrencyHandler).Methods("GET")
	v1.HandleFunc("/menus/{menuId:[0-9]+}", app.updateCurrencyHandler).Methods("PUT")
	v1.HandleFunc("/menus/{menuId:[0-9]+}", app.deleteCurrencyHandler).Methods("DELETE")

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
