package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (app *application) addDefaultData(td *templateData, c *gin.Context) *templateData {
	if td == nil {
		td = &templateData{}
	}

	year := time.Now().Year()

	td.CurrentYear = year

	s := sessions.Default(c)
	flash, ok := s.Get("flash").(string)
	if ok {
		s.Delete("flash")
		s.Save()
	}
	td.Flash = flash

	return td
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(c *gin.Context, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(c.Writer, fmt.Errorf("The template %s does not exist", name))
		return
	}

	buf := new(bytes.Buffer)

	err := ts.Execute(buf, app.addDefaultData(td, c))
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	buf.WriteTo(c.Writer)
}
