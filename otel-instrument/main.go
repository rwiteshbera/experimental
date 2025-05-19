package main

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	headers := map[string]string{
		"Authorization": "<>",
	}

	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint("localhost:4317"),
		otlpmetricgrpc.WithHeaders(headers),
		otlpmetricgrpc.WithDialOption(grpc.WithBlock()),
		otlpmetricgrpc.WithTemporalitySelector(func(kind sdkmetric.InstrumentKind) metricdata.Temporality {
			return metricdata.DeltaTemporality
		}),
	)
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}
	reader := sdkmetric.NewPeriodicReader(exporter,
		sdkmetric.WithInterval(5*time.Second),
	)

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(reader),
	)
	otel.SetMeterProvider(provider)

	meter := otel.Meter("app.example/metrics")

	histogram, err := meter.Float64Histogram("http_request_duration_seconds")
	if err != nil {
		log.Fatalf("failed to create histogram: %v", err)
	}

	histogram.Record(ctx, 2)
	histogram.Record(ctx, 3)
	histogram.Record(ctx, 8)
	time.Sleep(5 * time.Second)
	histogram.Record(ctx, 4)
	histogram.Record(ctx, 6)
	histogram.Record(ctx, 8)
	time.Sleep(5 * time.Second)
}
