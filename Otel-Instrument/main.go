package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	rand.Seed(time.Now().UnixNano())

	// Create OTLP gRPC exporter without headers
	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint("localhost:4317"),
		otlpmetricgrpc.WithDialOption(grpc.WithBlock()),
		otlpmetricgrpc.WithTemporalitySelector(func(kind sdkmetric.InstrumentKind) metricdata.Temporality {
			return metricdata.DeltaTemporality
		}),
	)
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}

	// Periodic reader every 5 seconds
	reader := sdkmetric.NewPeriodicReader(exporter,
		sdkmetric.WithInterval(5*time.Second),
	)

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(reader),
	)
	otel.SetMeterProvider(provider)

	meter := otel.Meter("app.example/metrics")

	// Create histogram instrument
	histogram, err := meter.Float64Histogram(
		"http_request_duration_seconds",
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(0.25, 0.5, 0.75, 1, 1.5, 2, 3),
	)
	if err != nil {
		log.Fatalf("failed to create histogram: %v", err)
	}

	// Continuously record realistic random request durations
	for {
		// Simulate 5–10 requests per interval
		for i := 0; i < rand.Intn(6)+5; i++ {
			var value float64
			r := rand.Float64()
			switch {
			case r < 0.7:
				// Fast requests: 0.1–0.5s
				value = 0.1 + rand.Float64()*0.4
			case r < 0.95:
				// Medium requests: 0.5–1s
				value = 0.5 + rand.Float64()*0.5
			default:
				// Slow requests: 1–3s
				value = 1.0 + rand.Float64()*2.0
			}

			fmt.Printf("Recording request duration: %.3f seconds\n", value)
			histogram.Record(ctx, value)
		}

		time.Sleep(5 * time.Second)
	}
}
