package main

import (
	"errors"
	"net/http"

	"github.com/nimaposhtiban/greenlight/internal/data"
	"github.com/nimaposhtiban/greenlight/internal/validator"
)

// @Summary Register a new user
// @Description Registers a new user with the provided info
// @BasePath /
// @Tags users
// @Accept json
// @Produce json
// @Param request body registerUserRequest true "Request body to register a user"
// @Success 201 "Created"
// @Failure 400 "Bad Request"
// @Failure 422 "Failed Model Validation"
// @Failure 500 "Internal Server Error"
// @Router /v1/users [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input registerUserRequest

	err := app.readJson(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)

		default:
			app.serverErrorResponse(w, r, err)
		}
	}
	err = app.writeJson(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
