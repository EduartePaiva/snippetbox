package main

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(app.recoverPanic, app.logRequest, secureHeaders)

	r.Group(func(r chi.Router) {
		r.Use(app.session.Enable, app.noSurf, app.authenticate)

		r.Get("/", app.home)
		r.With(app.requireAuthenticateUser).Get("/snippet/create", app.createSnippetForm)
		r.With(app.requireAuthenticateUser).Post("/snippet/create", app.createSnippet)
		r.Get("/snippet/{id}", app.showSnippet)
		r.Get("/user/signup", app.signupUserForm)
		r.Post("/user/signup", app.signupUser)
		r.Get("/user/login", app.loginUserForm)
		r.Post("/user/login", app.loginUser)
		r.With(app.requireAuthenticateUser).Post("/user/logout", app.logoutUser)

	})

	filesDir := http.Dir("./ui/static")
	FileServer(r, "/static/", filesDir)

	return r
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
