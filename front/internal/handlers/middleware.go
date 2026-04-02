package handlers

import (
	"net/http"

	"bibliotheque-front/internal/session"

	"github.com/gin-gonic/gin"
)

// RequireAuth redirects to /login if the user is not connected.
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := session.Get(c); !ok {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireBibliothecaire redirects to / if the connected user is not a bibliothécaire.
func RequireBibliothecaire() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := session.Get(c)
		if !ok || !user.IsBibliothecaire {
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}
		c.Next()
	}
}
