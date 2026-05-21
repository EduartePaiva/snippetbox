package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) routes() http.Handler {
	router := gin.New()

	router.Use(app.recoverPanic(), app.logRequest(), secureHeaders())

	router.GET("/", ginHandleFuncAdapter(app.home))
	router.POST("/snippet/create", ginHandleFuncAdapter(app.createSnippet))
	router.GET("/snippet/:id/", ginHandleFuncAdapter(app.showSnippet))

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
