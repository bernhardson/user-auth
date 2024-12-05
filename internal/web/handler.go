package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/bernhardson/stub/internal/models"
	"github.com/bernhardson/stub/internal/validator"
	"github.com/julienschmidt/httprouter"
)

func (app *Application) getUser(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	user, err := app.UserRepo.Get(int64(id))

	if err != nil {
		if user == nil {
			app.notFound(w)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (app *Application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	// Parse JSON into User struct
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	//validate inputs
	v := validator.Validator{}
	v.CheckField(validator.MinChars(user.Username, 3), "username", "User name must have at least three characters.")
	v.CheckField(validator.Matches(user.Email, validator.EmailRX), "email", "Entered email adress is not valid.")
	v.CheckField(validator.MinChars(user.Password, 8), "password", "Password must have at least 8 characters.")

	// Return input errors if any
	if !v.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": v.FieldErrors,
		})
		return
	}

	// Insert the user into the database
	_, err = app.UserRepo.Insert(user.Username, user.Email, user.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"errors": map[string]string{"email": "Email address already exists."},
			})
			return
		} else {
			app.serverError(w, err)
		}
	}

	// Send success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

type userLoginPost struct {
	Email    string
	Password string
}

func (app *Application) userLoginPost(w http.ResponseWriter, r *http.Request) {

	var body userLoginPost
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	v := validator.Validator{}

	v.CheckField(validator.NotBlank(body.Email), "email", "This field cannot be blank")
	v.CheckField(validator.Matches(body.Email, validator.EmailRX), "email", "Entered email adress is not valid")
	v.CheckField(validator.NotBlank(body.Password), "password", "This field cannot be blank")

	if !v.Valid() {
		app.clientError(w, http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": v.FieldErrors,
		})
		return
	}

	id, err := app.UserRepo.Authenticate(body.Email, body.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			v.AddNonFieldError("Email or password is incorrect")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"errors": v.NonFieldErrors,
			})
			return

		} else {
			app.serverError(w, err)
		}
		return
	}

	err = app.SessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.SessionManager.Put(r.Context(), "authenticatedUserID", id)
	// Send success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User logged in"})
}
func (app *Application) userLogoutPost(w http.ResponseWriter, r *http.Request) { // Use the RenewToken() method on the current session to change the session // ID again.

	err := app.SessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.SessionManager.Remove(r.Context(), "authenticatedUserID")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
