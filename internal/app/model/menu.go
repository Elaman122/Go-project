package model

import (
    "context"
    "database/sql"
    "log"
    "time"
)

type Menu struct {
    ID        int       `json:"id"`
    Code      string    `json:"code"`
    Rate      float64   `json:"rate"`
    Timestamp time.Time `json:"timestamp"`
}

type MenuModel struct {
    db        *sql.DB
    ErrorLog  *log.Logger
    InfoLog   *log.Logger
}

func (m MenuModel) Insert(menu *Menu) error {
    query := `
        INSERT INTO menu (code, rate, timestamp) 
        VALUES ($1, $2, $3) 
        RETURNING id
        `
    args := []interface{}{menu.Code, menu.Rate, menu.Timestamp}
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    return m.db.QueryRowContext(ctx, query, args...).Scan(&menu.ID)
}

func (m MenuModel) Get(id int) (*Menu, error) {
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

func (m MenuModel) Update(menu *Menu) error {
    query := `
        UPDATE menu
        SET code = $1, rate = $2, timestamp = $3
        WHERE id = $4
        `
    args := []interface{}{menu.Code, menu.Rate, menu.Timestamp, menu.ID}
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    _, err := m.db.ExecContext(ctx, query, args...)
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
