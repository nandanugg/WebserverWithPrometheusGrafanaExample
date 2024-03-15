package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	helloRequestHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "hello_request",
		Help:    "Histogram of the /hello request duration.",
		Buckets: prometheus.LinearBuckets(1, 1, 10), // Adjust bucket sizes as needed
	}, []string{"path", "method", "status"})
)

func main() {
	e := echo.New()

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	NewRoute(e, "/", "GET", helloHandler)
	NewRoute(e, "/:userId", "GET", helloHandler)

	e.Start(":8080")
}

func NewRoute(app *echo.Echo, path string, method string, handler echo.HandlerFunc) {
	app.Add(method, path, wrapHandlerWithMetrics(path, method, handler))
}

func wrapHandlerWithMetrics(path string, method string, handler echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		startTime := time.Now()

		// Execute the actual handler and catch any errors
		err := handler(c)

		// Regardless of whether an error occurred, record the metrics
		duration := time.Since(startTime).Seconds()

		helloRequestHistogram.WithLabelValues(path, method).Observe(duration)
		return err
	}
}

func helloHandler(c echo.Context) error {
	// Simulating processing time
	randomSeconds := rand.Intn(11)
	time.Sleep(time.Duration(randomSeconds) * time.Second)

	return c.String(http.StatusOK, "Hello, World! "+c.Param("userId"))
}
