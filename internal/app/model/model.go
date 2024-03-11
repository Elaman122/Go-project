package model

import (
    "database/sql"
    "log"
    "os"
)

type Models struct {
    Currency  CurrencyModel
    Menu      MenuModel
}

func NewModels(db *sql.DB) Models {
    infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
    errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
    return Models{
        Currency: CurrencyModel{
            db:       db,
            InfoLog:  infoLog,
            ErrorLog: errorLog,
        },
        Menu: MenuModel{
            db:       db,
            InfoLog:  infoLog,
            ErrorLog: errorLog,
        },
    }
}
