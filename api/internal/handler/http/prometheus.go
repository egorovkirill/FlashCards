package http

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/shirou/gopsutil/v3/cpu"
	"log"
	"strconv"
	"time"
)

var (
	CpuUsage = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cpu_usage",
		Help: "CPU usage percentage",
	}, []string{"core"})

	RequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "api_requests_total",
		Help: "Total number of API requests",
	}, []string{"method", "endpoint", "http_status"})

	RequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "api_request_duration_seconds",
		Help:    "API request duration distribution",
		Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2.5, 5},
	}, []string{"method", "endpoint", "http_status"})
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()

		duration := time.Since(startTime).Seconds()

		RequestsTotal.WithLabelValues(c.Request.Method, c.Request.URL.Path, strconv.Itoa(c.Writer.Status())).Inc()
		RequestDuration.WithLabelValues(c.Request.Method, c.Request.URL.Path, strconv.Itoa(c.Writer.Status())).Observe(duration)
	}
}

func UpdateCPUMetrics() {
	for {
		cpuUsageValues, err := cpu.Percent(0, false)
		if err != nil {
			log.Println(err)
			continue
		}

		for i, value := range cpuUsageValues {
			CpuUsage.WithLabelValues(strconv.Itoa(i)).Set(value)
		}

		time.Sleep(5 * time.Second)
	}
}
