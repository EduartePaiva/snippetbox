package main

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (app *application) routes() http.Handler {
	router := gin.New()

	router.Use(app.recoverPanic(), app.logRequest(), secureHeaders())

	dynamicRoutes := router.Group("")
	dynamicRoutes.Use(sessions.Sessions("session", app.store))
	{
		dynamicRoutes.GET("/", ginHandleFuncAdapter(app.home))
		dynamicRoutes.GET("/snippet/create", ginHandleFuncAdapter(app.createSnippetForm))
		dynamicRoutes.POST("/snippet/create", app.createSnippet)
		dynamicRoutes.GET("/snippet/:id/", app.showSnippet)
	}

	router.Static("/static/", "./ui/static")

	return router
}

func ginHandleFuncAdapter(f http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, param := range c.Params {
			c.Request.SetPathValue(param.Key, param.Value)
		}

		f(c.Writer, c.Request)
	}
}
