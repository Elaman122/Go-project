package model

import (
    "time"
	"database/sql"
	"errors"
	"log"
)

// Currency ...
type Currency struct {
    ID        int       `json:"id"`
    Code      string    `json:"code"`
    Rate      float64   `json:"rate"`
    Timestamp time.Time `json:"timestamp"`
}

//ss
type CurrencyModel struct {
	db             *sql.DB
	ErrorLog 	   *log.Logger
	InfoLog  	   *log.Logger
}

var currency = []Currency {
	{ID: 1, Code: "USD", Rate: 1.0, Timestamp: time.Now()},
    {ID: 2, Code: "EUR", Rate: 0.84, Timestamp: time.Now()},
	{ID: 3, Code: "GBP", Rate: 1.23, Timestamp: time.Now()},    
    {ID: 4, Code: "RUB", Rate: 0.011, Timestamp: time.Now()},  
    {ID: 5, Code: "KZT", Rate: 0.0023, Timestamp: time.Now()},
}

func getCurrency() []Currency {
	return currency
}

func getCurrencies(id int) (*Currency, error) {
    for _, r := range currency {
        if r.ID == id {
            return &r, nil
        }
    }

    return nil, errors.New("Not found")
}