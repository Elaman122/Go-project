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
		// Чтение данных создания меню из тела запроса
		var input model.Menu
		err := app.readJSON(w, r, &input)
		if err != nil {
			app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
	
		// Проверка разрешения CREATE для пользователя
		user := app.contextGetUser(r)
		requiredPermissionID := 2 // Установите значение, соответствующее разрешению для создания
		hasPermission, err := app.models.Permissions.CheckPermission(user.ID, requiredPermissionID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		if !hasPermission {
			app.notPermittedResponse(w, r)
			return
		}
	
		// Вставка нового меню
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
		// Извлечение ID элемента меню из URL
		vars := mux.Vars(r)
		param := vars["menuId"]
	
		id, err := strconv.Atoi(param)
		if err != nil || id < 1 {
			app.respondWithError(w, http.StatusBadRequest, "Invalid currency ID")
			return
		}
	
		// Получение данных меню для обновления
		currency, err := app.models.Menu.Get(id)
		if err != nil {
			app.respondWithError(w, http.StatusNotFound, "404 Not Found")
			return
		}
	
		// Чтение данных обновления из тела запроса
		var input model.Menu
		err = app.readJSON(w, r, &input)
		if err != nil {
			app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
	
		// Проверка разрешения UPDATE для пользователя
		user := app.contextGetUser(r)
		requiredPermissionID := 2 // Установите значение, соответствующее разрешению для обновления
		hasPermission, err := app.models.Permissions.CheckPermission(user.ID, requiredPermissionID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		if !hasPermission {
			app.notPermittedResponse(w, r)
			return
		}
	
		// Обновление данных меню
		currency.Code = input.Code
		currency.Rate = input.Rate
		currency.Timestamp = input.Timestamp
		currency.CurrencyCode = input.CurrencyCode
	
		err = app.models.Menu.Update(currency)
		if err != nil {
			app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
			return
		}
	
		app.respondWithJSON(w, http.StatusOK, currency)
	}

	func (app *application) deleteCurrencyHandler(w http.ResponseWriter, r *http.Request) {
		// Извлечь ID элемента меню из URL
		vars := mux.Vars(r)
		param := vars["menuId"]
		id, err := strconv.Atoi(param)
		if err != nil || id < 1 {
			app.respondWithError(w, http.StatusBadRequest, "Invalid menu ID")
			return
		}
	
		// Проверить разрешение DELETE для пользователя
		user := app.contextGetUser(r)
		requiredPermissionID := 2 // Замените это значением, которое соответствует разрешению для удаления
		hasPermission, err := app.models.Permissions.CheckPermission(user.ID, requiredPermissionID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		if !hasPermission {
			app.notPermittedResponse(w, r)
			return
		}
	
		// Выполнить удаление элемента меню
		err = app.models.Menu.Delete(id)
		if err != nil {
			app.respondWithError(w, http.StatusInternalServerError, "Failed to delete menu item")
			return
		}
	
		app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
	}

	
	func (app *application) getAllMenuHandler(w http.ResponseWriter, r *http.Request) {
		// Извлечение параметров из URL запроса
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
		sort := r.URL.Query().Get("sort")
		code := r.URL.Query().Get("code")
		from, _ := strconv.Atoi(r.URL.Query().Get("from"))
		to, _ := strconv.Atoi(r.URL.Query().Get("to"))

	
		// Создание Filters объекта на основе извлеченных параметров
		filters := model.Filters{
			Page:         page,
			PageSize:     pageSize,
			Sort:         sort,
			SortSafeList: []string{"rate", "code", "timestamp"}, // Правильное заполнение SortSafeList
		}

		// Вызов метода GetAll вашей MenuModel с фильтрами
		menu, metadata, err := app.models.Menu.GetAll(code, from, to, filters)
		if err != nil {
			// Обработка ошибки, если запрос не удалось выполнить
			app.respondWithError(w, http.StatusInternalServerError, "Failed to fetch menus")
			return
		}

		// Отправка ответа с данными меню и метаданными пагинации
		response := map[string]interface{}{
			"menus":    menu,
			"metadata": metadata,
		}
		app.respondWithJSON(w, http.StatusOK, response)
	}

