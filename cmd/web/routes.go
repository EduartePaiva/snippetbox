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
		dynamicRoutes.GET("/", app.home)
		dynamicRoutes.GET("/snippet/create", app.createSnippetForm)
		dynamicRoutes.POST("/snippet/create", app.createSnippet)
		dynamicRoutes.GET("/snippet/:id/", app.showSnippet)
		dynamicRoutes.GET("/user/signup", app.signupUserForm)
		dynamicRoutes.POST("/user/signup", app.signupUser)
		dynamicRoutes.GET("/user/login", app.loginUserForm)
		dynamicRoutes.POST("/user/login", app.loginUser)
		dynamicRoutes.POST("/user/logout", app.logoutUser)
	}

	router.Static("/static/", "./ui/static")

	return router
}
