package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"guthub.com/eduartepaiva/snippetbox/pkg/forms"
	"guthub.com/eduartepaiva/snippetbox/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.html", &templateData{Snippets: s})
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	numId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || numId < 1 {
		app.notFound(w)
		return
	}
	s, err := app.snippets.Get(numId)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	}
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "show.page.html", &templateData{Snippet: s})
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.html", &templateData{Form: forms.New(nil)})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(w, r, "create.page.html", &templateData{
			Form: form,
		})
		return
	}

	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Snippet successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.html", &templateData{Form: forms.New(nil)})
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MinLength("password", 10)
	form.MatchesPattern("email", forms.EmailRX)

	if !form.Valid() {
		app.render(w, r, "signup.page.html", &templateData{Form: form})
		return
	}

	_, err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if errors.Is(err, models.ErrDuplicateEmail) {
		form.Errors.Add("email", "Address is already in use")
		app.render(w, r, "signup.page.html", &templateData{Form: form})
		return
	}

	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.html", &templateData{Form: forms.New(nil)})
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.MatchesPattern("email", forms.EmailRX)

	if !form.Valid() {
		form.Errors.Add("generic", "invalid email or password")
		app.render(w, r, "login.page.html", &templateData{Form: form})
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	id, err := app.users.Authenticate(email, password)
	if err == models.ErrInvalidCredentials {
		form.Errors.Add("generic", "Email or Password is incorrect")
		app.render(w, r, "login.page.html", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "userID", id)
	app.session.Put(r, "flash", "You logged in!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "userID")
	app.session.Put(r, "flash", "You've been logged out!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
