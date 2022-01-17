package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/justinas/nosurf"
)

func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

func AddRequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := uuid.NewV4()
		c.Writer.Header().Set("X-Request-Id", id.String())
		c.Next()
	}
}
