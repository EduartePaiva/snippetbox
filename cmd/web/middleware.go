package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func secureHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("X-Frame-Options", "deny")
		c.Writer.Header().Add("X-XSS-Protection", "1; mode=block")

		c.Next()
	}
}

func (app *application) logRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		app.infoLog.Printf(
			"%s - %s %s %s",
			c.Request.RemoteAddr,
			c.Request.Proto,
			c.Request.Method,
			c.Request.URL.RequestURI(),
		)

		c.Next()
	}
}

func (app *application) recoverPanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.Writer.Header().Set("Connection", "close")
				app.serverError(c.Writer, fmt.Errorf("%s", err))
			}
		}()

		c.Next()
	}
}
