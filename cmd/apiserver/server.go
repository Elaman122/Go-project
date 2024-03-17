package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Elaman122/Go-project/internal/app/model"
	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type server struct {
	router       *mux.Router
	logger       *logrus.Logger
	sessionStore sessions.Store
	application  *application
}

func (app *application) respondWithError(w http.ResponseWriter, code int, message string) {
	app.respondWithJSON(w, code, map[string]string{"error": message})
}

func (app *application) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (app *application) createCurrencyHandler(w http.ResponseWriter, r *http.Request) {
	var input model.Menu

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = app.models.Menu.Insert(&input)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, input)
}

func (app *application) getCurrencyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["menuId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid currency ID")
		return
	}

	currency, err := app.models.Menu.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, currency)
}

func (app *application) updateCurrencyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["menuId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid currency ID")
		return
	}

	currency, err := app.models.Menu.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input model.Menu

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	currency.Code = input.Code
	currency.Rate = input.Rate
	currency.Timestamp = input.Timestamp

	err = app.models.Menu.Update(currency)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, currency)
}

func (app *application) deleteCurrencyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["menuId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid currency ID")
		return
	}

	err = app.models.Menu.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		return err
	}

	return nil
}
