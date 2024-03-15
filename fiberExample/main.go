package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Define the histogram metric.
var (
	helloRequestHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "hello_request",
		Help:    "Histogram of the /hello request duration.",
		Buckets: prometheus.LinearBuckets(1, 1, 10), // Adjust bucket sizes as needed
	}, []string{"path", "method", "status"})
)

func main() {
	app := fiber.New()

	// Use recover middleware to prevent crashes
	app.Use(recover.New())

	// Register the Prometheus metrics handler
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	// Simplified route creation with NewRoute function
	NewRoute(app, "/", "GET", helloHandler)
	NewRoute(app, "/:userId", "GET", helloHandler)

	// Start the server
	fmt.Println("WEBSERVER IS RUNNING ON :8080")
	log.Fatal(app.Listen(":8080"))
}

func NewRoute(app *fiber.App, path string, method string, handler fiber.Handler) {
	app.Add(method, path, wrapHandlerWithMetrics(path, method, handler))
}

func wrapHandlerWithMetrics(path string, method string, handler fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		startTime := time.Now()

		// Execute the actual handler and catch any errors
		err := handler(c)

		// Regardless of whether an error occurred, record the metrics
		duration := time.Since(startTime).Seconds()
		statusCode := fmt.Sprintf("%d", c.Response().StatusCode())
		if err != nil {
			if c.Response().StatusCode() == fiber.StatusOK { // Default status code
				statusCode = "500" // Assume internal server error if not set
			}
			c.Status(fiber.StatusInternalServerError).SendString(err.Error()) // Ensure the response reflects the error
		}

		helloRequestHistogram.WithLabelValues(path, method, statusCode).Observe(duration)
		return err
	}
}

func helloHandler(c *fiber.Ctx) error {
	// Simulating processing time
	randomSeconds := rand.Intn(11)
	time.Sleep(time.Duration(randomSeconds) * time.Second)

	return c.SendString("hello world " + c.Params("userId", ""))
}
