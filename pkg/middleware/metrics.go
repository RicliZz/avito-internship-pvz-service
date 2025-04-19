package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

var CountAllRequests = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "total_requests",
	Help: "Total number of requests",
})

var buckets = []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}

var MeasureResponseDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "measure_response_duration_seconds",
	Help:    "Measure response duration in seconds",
	Buckets: buckets,
}, []string{"route", "method", "status_code"},
)

func ResponseDurationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start).Seconds()
		MeasureResponseDuration.WithLabelValues(c.FullPath(), c.Request.Method,
			strconv.Itoa(c.Writer.Status())).Observe(duration)
	}
}

func RequestCounterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		CountAllRequests.Inc()
		c.Next()
	}
}
