package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)

	r.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	r.HandleFunc("/api/v1/healthcheck", app.healthcheckHandler).Methods("GET")

	v1 := r.PathPrefix("/api/v1").Subrouter()

	// Обработчики для создания, получения, обновления и удаления элементов меню
	v1.HandleFunc("/menus", app.createCurrencyHandler).Methods("POST")
	v1.HandleFunc("/menus/{menuId:[0-9]+}", app.getCurrencyHandler).Methods("GET")
	v1.HandleFunc("/menus/{menuId:[0-9]+}", app.updateCurrencyHandler).Methods("PUT")

	// Route that requires DELETE permission for deleting menu items
	v1.HandleFunc("/menus/{menuId:[0-9]+}", app.deleteCurrencyHandler).Methods("DELETE")

	// Обработчик для получения списка меню с поддержкой пагинации, сортировки и фильтрации
	v1.HandleFunc("/menufor", app.getAllMenuHandler).Methods("GET")

	users1 := r.PathPrefix("/api/v1").Subrouter()
	// User handlers with Authentication
	users1.HandleFunc("/users", app.registerUserHandler).Methods("POST")
	users1.HandleFunc("/users/activated", app.activateUserHandler).Methods("PUT")
	users1.HandleFunc("/users/login", app.createAuthenticationTokenHandler).Methods("POST")

	return app.authenticate(r)
}