package filler

import (
	"time"

	model "github.com/Elaman122/Go-project/internal/app/model"
)

func PopulateDatabase(models model.Models) error {
	menus := []model.Menu{
		{ID: 1, Code: "Menu1", Rate: 10.99, Timestamp: time.Now(), CurrencyCode: "USD"},
		{ID: 2, Code: "Menu2", Rate: 15.99, Timestamp: time.Now(), CurrencyCode: "EUR"},
		// Add more menus as needed
	}

	for _, menu := range menus {
		models.Menu.Insert(&menu)
	}
	// TODO: Implement restaurants population
	// TODO: Implement the relationship between restaurants and menus
	return nil
}