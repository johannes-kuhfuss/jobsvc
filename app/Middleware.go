package app

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/misc"
)

func AddRequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := uuid.NewV4()
		c.Writer.Header().Set("X-Request-Id", id.String())
		c.Next()
	}
}

func validateApiKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := strings.TrimSpace(c.Request.Header.Get("Authorization"))
		split := strings.Split(apiKey, "Bearer ")
		if len(split) == 2 {
			if misc.SliceContainsString(cfg.Misc.ApiKeys, split[1]) {
				return
			} else {
				err := api_error.NewUnauthenticatedError("Could not verify API key")
				c.AbortWithStatusJSON(err.StatusCode(), err)
			}
		} else {
			err := api_error.NewUnauthenticatedError("Could not verify API key")
			c.AbortWithStatusJSON(err.StatusCode(), err)
		}
	}
}
