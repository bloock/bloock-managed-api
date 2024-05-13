package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of get requests.",
	},
	[]string{"path"},
)

var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "response_status",
		Help: "Status of HTTP response",
	},
	[]string{"status"},
)

/*var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_response_time_seconds",
	Help: "Duration of HTTP requests.",
}, []string{"path"})*/

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.String()

		//timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))

		c.Next()

		responseStatus.WithLabelValues(strconv.Itoa(c.Writer.Status())).Inc()
		totalRequests.WithLabelValues(path).Inc()

		//timer.ObserveDuration()
	}
}

func NewMetricsMiddleware() error {
	_ = prometheus.Register(totalRequests)
	_ = prometheus.Register(responseStatus)
	//_ = prometheus.Register(httpDuration)

	return nil
}
