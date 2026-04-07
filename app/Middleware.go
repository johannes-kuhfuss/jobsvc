package app

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/misc"
	"github.com/prometheus/client_golang/prometheus"
)

func AddRequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := uuid.NewV4()
		c.Writer.Header().Set("X-Request-Id", id.String())
		c.Next()
	}
}

func validateAuth() gin.HandlerFunc {
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
		c.Next()
	}
}

func prometheusMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		c.Next()
		statusCode := c.Writer.Status()
		responseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
		totalRequests.WithLabelValues(path).Inc()
		timer.ObserveDuration()
	}
}
