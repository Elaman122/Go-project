package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/Elaman122/Go-project/internal/app/model"
	"github.com/Elaman122/Go-project/internal/app/validator"
)

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Извлечь email и пароль из тела запроса.
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Проверить корректность email и пароля, переданных клиентом.
	v := validator.New()
	model.ValidateEmail(v, input.Email)
	model.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Найти запись пользователя по email. Если пользователь не найден, отправить ошибку 401.
	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Проверить соответствие переданного пароля фактическому паролю пользователя.
	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Если пароли не совпадают, отправить ошибку 401 и завершить обработку запроса.
	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	// Иначе, если пароль верный, сгенерировать новый токен с сроком действия 24 часа
	// и сферой "аутентификация".
	token, err := app.models.Tokens.New(user.ID, 24*time.Hour, model.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Кодировать токен в формат JSON и отправить его в ответе с кодом 201 Created.
	err = app.writeJSON(w, http.StatusCreated, envelope{"authentication_token": token}, nil)
	if err != nil {
    	app.serverErrorResponse(w, r, err)
    	return
	}
}