package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"guthub.com/eduartepaiva/snippetbox/pkg/forms"
	"guthub.com/eduartepaiva/snippetbox/pkg/models"
)

func (app *application) home(c *gin.Context) {
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	app.render(c, "home.page.html", &templateData{Snippets: s})
}

func (app *application) showSnippet(c *gin.Context) {
	numId, err := strconv.Atoi(c.Param("id"))
	if err != nil || numId < 1 {
		app.notFound(c.Writer)
		return
	}
	s, err := app.snippets.Get(numId)
	if err == models.ErrNoRecord {
		app.notFound(c.Writer)
		return
	}
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	app.render(c, "show.page.html", &templateData{Snippet: s})
}

func (app *application) createSnippetForm(c *gin.Context) {
	app.render(c, "create.page.html", &templateData{Form: forms.New(nil)})
}

func (app *application) createSnippet(c *gin.Context) {
	session := sessions.Default(c)
	err := c.Request.ParseForm()
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}
	form := forms.New(c.Request.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(c, "create.page.html", &templateData{
			Form: form,
		})
		return
	}

	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	session.Set("flash", "Snippet successfully created!")
	session.Save()
	http.Redirect(c.Writer, c.Request, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) signupUserForm(c *gin.Context) {
	app.render(c, "signup.page.html", &templateData{Form: forms.New(nil)})
}

func (app *application) signupUser(c *gin.Context) {
	fmt.Fprintln(c.Writer, "create a new user...")
}

func (app *application) loginUserForm(c *gin.Context) {
	fmt.Fprintln(c.Writer, "display the user login form...")
}

func (app *application) loginUser(c *gin.Context) {
	fmt.Fprintln(c.Writer, "authenticate and login the user...")
}

func (app *application) logoutUser(c *gin.Context) {
	fmt.Fprintln(c.Writer, "logout the user...")
}
