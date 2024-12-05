package web

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/bernhardson/stub/internal/repo"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/rs/zerolog"
)

type Application struct {
	UserRepo       repo.UserRepository
	Logger         *zerolog.Logger
	SessionManager *scs.SessionManager
}

func (app *Application) Routes() http.Handler {

	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	dynamic := alice.New(app.SessionManager.LoadAndSave, app.authenticate)
	// dynamicJson := dynamic.Append(RequireJSON)
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodPost, "/user/signup", standard.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodPost, "/user/logout", dynamic.ThenFunc(app.userLogoutPost))
	router.Handler(http.MethodGet, "/user/view/:id", protected.ThenFunc(app.getUser))
	router.Handler(http.MethodGet, "/ping", http.HandlerFunc(app.Ping))

	return standard.Then(router)
}
