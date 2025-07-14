package main

import (
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	requestDurations := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds.",
			Buckets: []float64{0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2, 2.5},
		},
		[]string{"method", "path"},
	)
	prometheus.MustRegister(requestDurations)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	for {
		value := realisticLatency()
		requestDurations.WithLabelValues("GET", "/foo").Observe(value)
		time.Sleep(10 * time.Millisecond)
	}
}

func realisticLatency() float64 {
	// Skew distribution more toward lower values (e.g., 0.01 to 0.3), but still with spikes
	mean := math.Log(0.05) // center around ~50ms
	stddev := 0.7          // wider spread for more variety

	raw := math.Exp(rand.NormFloat64()*stddev + mean)
	return math.Min(2.5, math.Max(0.001, raw))
}
