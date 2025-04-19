package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var CountAllRequests = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "total_requests",
	Help: "Total number of requests",
})

func RequestCounterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		CountAllRequests.Inc()
		c.Next()
	}
}
