package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
	"unicode/utf8"

	"github.com/Elaman122/Go-project/internal/app/validator"
)

type Menu struct {
	ID          int       `json:"id"`
	Code        string    `json:"code"`
	Rate        float64   `json:"rate"`
	Timestamp   time.Time `json:"timestamp"`
	CurrencyCode string   `json:"currency_code"`
}

type MenuModel struct {
	db       *sql.DB
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}


func (m *MenuModel) GetAll(code string, from, to int, filters Filters) ([]*Menu, int, error) {
    query := `
		SELECT id, code, rate, timestamp
		FROM menu
		WHERE ($1 = '' OR LOWER(code) = LOWER($1))
		AND ($2 = 0 OR rate >= $2)
		AND ($3 = 0 OR rate <= $3)
		ORDER BY %s %s
		LIMIT $4 OFFSET $5
	`
    query = fmt.Sprintf(query, filters.sortColumn(), filters.sortDirection())

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    args := []interface{}{code, from, to, filters.PageSize, filters.Offset} // Исправлено на filters.Offset

    rows, err := m.db.QueryContext(ctx, query, args...)
    if err != nil {
        m.ErrorLog.Println("Error querying database:", err)
        return nil, 0, err
    }
    defer rows.Close()

    var menus []*Menu
    for rows.Next() {
        var menu Menu
        err := rows.Scan(&menu.ID, &menu.Code, &menu.Rate, &menu.Timestamp)
        if err != nil {
            m.ErrorLog.Println("Error scanning rows:", err)
            return nil, 0, err
        }
        menus = append(menus, &menu)
    }

    var totalRecords int
    err = m.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM menu").Scan(&totalRecords)
    if err != nil {
        m.ErrorLog.Println("Error getting total records:", err)
        return nil, 0, err
    }

    return menus, totalRecords, nil
}

func (m *MenuModel) GetPaginated(page int, pageSize int, sort string) ([]*Menu, error) {
    if sort != "code" && sort != "rate" && sort != "timestamp" {
        return nil, errors.New("Invalid sort parameter")
    }

    query := fmt.Sprintf("SELECT id, code, rate, timestamp FROM menu ORDER BY %s LIMIT $1 OFFSET $2", sort)

    rows, err := m.db.Query(query, pageSize, (page-1)*pageSize)
    if err != nil {
        m.ErrorLog.Println("Error querying database for pagination:", err)
        return nil, err
    }
    defer rows.Close()

    var menus []*Menu
    for rows.Next() {
        var menu Menu
        if err := rows.Scan(&menu.ID, &menu.Code, &menu.Rate, &menu.Timestamp); err != nil {
            m.ErrorLog.Println("Error scanning rows for pagination:", err)
            return nil, err
        }
        menus = append(menus, &menu)
    }

    if err := rows.Err(); err != nil {
        m.ErrorLog.Println("Error during pagination rows iteration:", err)
        return nil, err
    }

    return menus, nil
}


func (m *MenuModel) Insert(menu *Menu) error {
	query := `
		INSERT INTO menu (code, rate, timestamp) 
		VALUES ($1, $2, $3) 
		RETURNING id
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.db.QueryRowContext(ctx, query, menu.Code, menu.Rate, menu.Timestamp).Scan(&menu.ID)
}

func (m *MenuModel) Get(id int) (*Menu, error) {
	query := `
		SELECT id, code, rate, timestamp
		FROM menu
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var menu Menu
	err := m.db.QueryRowContext(ctx, query, id).Scan(&menu.ID, &menu.Code, &menu.Rate, &menu.Timestamp)
	if err != nil {
		return nil, err
	}

	return &menu, nil
}

func (m *MenuModel) Update(menu *Menu) error {
	query := `
		UPDATE menu
		SET code = $1, rate = $2, timestamp = $3
		WHERE id = $4
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.db.ExecContext(ctx, query, menu.Code, menu.Rate, menu.Timestamp, menu.ID)
	return err
}

func (m *MenuModel) Delete(id int) error {
	query := `
		DELETE FROM menu
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.db.ExecContext(ctx, query, id)
	return err
}


func ValidateMenu(v *validator.Validator, menu *Menu) {
	v.Check(menu.Code != "", "code", "must be provided")
	v.Check(utf8.RuneCountInString(menu.Code) <= 100, "code", "must be between 1 and 100 characters long")
	v.Check(menu.Rate <= 100000.00, "rate", "must not be more than 100000.00")
}
