package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	metricspb "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	otlpmetrics "go.opentelemetry.io/proto/otlp/metrics/v1"
	resourcepb "go.opentelemetry.io/proto/otlp/resource/v1"
)

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	conn, err := grpc.Dial("localhost:4317", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer conn.Close()

	// add metadata to the request
	md := metadata.New(map[string]string{
		"Authorization": "<your_authorization_token>",
	})

	client := metricspb.NewMetricsServiceClient(conn)

	// Send each metric type with a 1-second delay
	metricTypes := []string{
		"monotonic_cumulative_sum",
		// "non_monotonic_sum",
		// "delta_sum",
		// "gauge",
		// "cumulative_histogram",
		// "delta_histogram",
	}

	for _, metricType := range metricTypes {
		req := buildMetricsRequestForType(metricType)
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		resp, err := client.Export(ctx, req)
		if err != nil {
			log.Printf("gRPC Export failed for %s: %v", metricType, err)
			continue
		}

		log.Printf("Export response for %s: %+v", metricType, resp)

		// Wait for 1 second before sending next metric
		time.Sleep(1 * time.Second)
	}
}

func buildMetricsRequestForType(metricType string) *metricspb.ExportMetricsServiceRequest {
	now := uint64(time.Now().UnixNano())
	timestamp10sAgo := uint64(time.Now().Add(-10 * time.Second).UnixNano())

	// Create base request with resource
	req := &metricspb.ExportMetricsServiceRequest{
		ResourceMetrics: []*otlpmetrics.ResourceMetrics{
			{
				Resource: &resourcepb.Resource{
					Attributes: []*commonpb.KeyValue{
						{Key: "service.name", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "grpc-metrics-demo"}}},
						{Key: "environment", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "development"}}},
					},
				},
				ScopeMetrics: []*otlpmetrics.ScopeMetrics{
					{
						Metrics: []*otlpmetrics.Metric{},
					},
				},
			},
		},
	}

	// Add specific metric based on type
	switch metricType {
	case "monotonic_cumulative_sum":
		req.ResourceMetrics[0].ScopeMetrics[0].Metrics = append(req.ResourceMetrics[0].ScopeMetrics[0].Metrics,
			&otlpmetrics.Metric{
				Name:        "KloudMate.requests.total",
				Description: "Total number of requests (monotonic cumulative counter)",
				Unit:        "requests",
				Data: &otlpmetrics.Metric_Sum{
					Sum: &otlpmetrics.Sum{
						DataPoints: []*otlpmetrics.NumberDataPoint{
							{
								TimeUnixNano:      now,
								StartTimeUnixNano: now - uint64(60*time.Second.Nanoseconds()),
								Value:             &otlpmetrics.NumberDataPoint_AsDouble{AsDouble: 1},
								Attributes: []*commonpb.KeyValue{
									{Key: "endpoint", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "/api/v1/users"}}},
									{Key: "method", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "GET"}}},
									{Key: "__name__", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "KloudMate.requests.total"}}},
									{Key: "job", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "grpc-test"}}},
									{Key: "defaultWorkspaceId", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "test-workspace"}}},
								},
							},
						},
						AggregationTemporality: otlpmetrics.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
						IsMonotonic:            true,
					},
				},
			})

	case "non_monotonic_sum":
		req.ResourceMetrics[0].ScopeMetrics[0].Metrics = append(req.ResourceMetrics[0].ScopeMetrics[0].Metrics,
			&otlpmetrics.Metric{
				Name:        "KloudMate.active.goroutines",
				Description: "Simulated number of active goroutines (can increase or decrease)",
				Unit:        "goroutines",
				Data: &otlpmetrics.Metric_Sum{
					Sum: &otlpmetrics.Sum{
						DataPoints: []*otlpmetrics.NumberDataPoint{
							{
								TimeUnixNano:      now,
								StartTimeUnixNano: timestamp10sAgo,
								Value:             &otlpmetrics.NumberDataPoint_AsDouble{AsDouble: float64(10 + rand.Intn(20) - 10)},
								Attributes: []*commonpb.KeyValue{
									{Key: "instance", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "grpc-node-1"}}},
									{Key: "__name__", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "KloudMate.active.goroutines"}}},
									{Key: "job", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "grpc-test"}}},
									{Key: "defaultWorkspaceId", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "test-workspace"}}},
								},
							},
						},
						AggregationTemporality: otlpmetrics.AggregationTemporality_AGGREGATION_TEMPORALITY_DELTA,
						IsMonotonic:            false,
					},
				},
			})

	case "delta_sum":
		req.ResourceMetrics[0].ScopeMetrics[0].Metrics = append(req.ResourceMetrics[0].ScopeMetrics[0].Metrics,
			&otlpmetrics.Metric{
				Name:        "KloudMate.requests.delta",
				Description: "Delta number of requests since last measurement",
				Unit:        "requests",
				Data: &otlpmetrics.Metric_Sum{
					Sum: &otlpmetrics.Sum{
						DataPoints: []*otlpmetrics.NumberDataPoint{
							{
								TimeUnixNano:      now,
								StartTimeUnixNano: timestamp10sAgo,
								Value:             &otlpmetrics.NumberDataPoint_AsDouble{AsDouble: 8},
								Attributes: []*commonpb.KeyValue{
									{Key: "endpoint", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "/api/v1/metrics"}}},
									{Key: "method", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "POST"}}},
									{Key: "__name__", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "KloudMate.requests.delta"}}},
									{Key: "job", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "grpc-test"}}},
									{Key: "defaultWorkspaceId", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "test-workspace"}}},
								},
							},
						},
						AggregationTemporality: otlpmetrics.AggregationTemporality_AGGREGATION_TEMPORALITY_DELTA,
						IsMonotonic:            true,
					},
				},
			})

	case "gauge":
		req.ResourceMetrics[0].ScopeMetrics[0].Metrics = append(req.ResourceMetrics[0].ScopeMetrics[0].Metrics,
			&otlpmetrics.Metric{
				Name:        "KloudMate.system.memory.usage",
				Description: "Current memory usage (gauge)",
				Unit:        "bytes",
				Data: &otlpmetrics.Metric_Gauge{
					Gauge: &otlpmetrics.Gauge{
						DataPoints: []*otlpmetrics.NumberDataPoint{
							{
								TimeUnixNano: now,
								Value:        &otlpmetrics.NumberDataPoint_AsDouble{AsDouble: 950},
								Attributes: []*commonpb.KeyValue{
									{Key: "host", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "test-host"}}},
									{Key: "__name__", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "KloudMate.system.memory.usage"}}},
									{Key: "job", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "grpc-test"}}},
									{Key: "defaultWorkspaceId", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "test-workspace"}}},
								},
							},
						},
					},
				},
			})

	case "cumulative_histogram":
		req.ResourceMetrics[0].ScopeMetrics[0].Metrics = append(req.ResourceMetrics[0].ScopeMetrics[0].Metrics,
			&otlpmetrics.Metric{
				Name:        "KloudMate.request.duration",
				Description: "Request duration distribution",
				Unit:        "ms",
				Data: &otlpmetrics.Metric_Histogram{
					Histogram: &otlpmetrics.Histogram{
						DataPoints: []*otlpmetrics.HistogramDataPoint{
							{
								TimeUnixNano:      now,
								StartTimeUnixNano: timestamp10sAgo,
								Count:             50,
								Sum:               func() *float64 { f := 1000.0; return &f }(),
								ExplicitBounds:    []float64{5, 10, 25, 50, 100, 250, 500, 1000},
								BucketCounts:      []uint64{5, 10, 15, 10, 5, 3, 2, 0, 0},
								Attributes: []*commonpb.KeyValue{
									{Key: "endpoint", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "/api/v1/data"}}},
									{Key: "__name__", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "KloudMate.request.duration"}}},
									{Key: "job", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "grpc-test"}}},
									{Key: "defaultWorkspaceId", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "test-workspace"}}},
								},
							},
						},
						AggregationTemporality: otlpmetrics.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
					},
				},
			})

	case "delta_histogram":
		req.ResourceMetrics[0].ScopeMetrics[0].Metrics = append(req.ResourceMetrics[0].ScopeMetrics[0].Metrics,
			&otlpmetrics.Metric{
				Name:        "KloudMate.request.duration.delta",
				Description: "Request duration distribution (delta)",
				Unit:        "ms",
				Data: &otlpmetrics.Metric_Histogram{
					Histogram: &otlpmetrics.Histogram{
						DataPoints: []*otlpmetrics.HistogramDataPoint{
							{
								TimeUnixNano:      now,
								StartTimeUnixNano: timestamp10sAgo,
								Count:             10,
								Sum:               func() *float64 { f := 200.0; return &f }(),
								ExplicitBounds:    []float64{5, 10, 25, 50, 100},
								BucketCounts:      []uint64{1, 2, 3, 2, 1, 1},
								Attributes: []*commonpb.KeyValue{
									{Key: "endpoint", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "/api/v1/status"}}},
									{Key: "__name__", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "KloudMate.request.duration.delta"}}},
									{Key: "job", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "grpc-test"}}},
									{Key: "defaultWorkspaceId", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "test-workspace"}}},
								},
							},
						},
						AggregationTemporality: otlpmetrics.AggregationTemporality_AGGREGATION_TEMPORALITY_DELTA,
					},
				},
			})
	}

	return req
}
