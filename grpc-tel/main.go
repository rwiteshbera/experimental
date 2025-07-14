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

	client := metricspb.NewMetricsServiceClient(conn)

	md := metadata.New(map[string]string{
		"Authorization": "sk_os3wltDw6M9jl3baCdGxhK8E",
	})

	// Send each metric type with a 1-second delay
	metricTypes := []string{
		// "monotonic_cumulative_sum",
		// "non_monotonic_sum",
		// "delta_sum",
		// "gauge",
		"cumulative_histogram",
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
						{Key: "service.name", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "grpc-metrics-demo-2"}}},
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
								Value:             &otlpmetrics.NumberDataPoint_AsDouble{AsDouble: 20},
								Attributes: []*commonpb.KeyValue{
									{Key: "endpoint", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "/api/v1/users"}}},
									{Key: "method", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "GET"}}},
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
								Value:             &otlpmetrics.NumberDataPoint_AsDouble{AsDouble: 35},
								Attributes: []*commonpb.KeyValue{
									{Key: "instance", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "grpc-node-1"}}},
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
								Value:             &otlpmetrics.NumberDataPoint_AsDouble{AsDouble: 100},
								Attributes: []*commonpb.KeyValue{
									{Key: "endpoint", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "/api/v1/metrics"}}},
									{Key: "method", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "POST"}}},
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
								Value:        &otlpmetrics.NumberDataPoint_AsDouble{AsDouble: 150},
								Attributes: []*commonpb.KeyValue{
									{Key: "host", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "test-host"}}},
								},
							},
						},
					},
				},
			})

	case "cumulative_histogram":
		bounds := []float64{50, 100, 150, 200, 250, 300}
		bucketCounts := []uint64{50, 40, 30, 20, 10, 5} // len = len(bounds) + 1

		// Compute total count
		var totalCount uint64
		for _, c := range bucketCounts {
			totalCount += c
		}

		// Estimate sum using bucket midpoints
		var totalSum float64
		for i, count := range bucketCounts {
			var midpoint float64
			if i < len(bounds) {
				lower := 0.0
				if i > 0 {
					lower = bounds[i-1]
				}
				upper := bounds[i]
				midpoint = (lower + upper) / 2
			} else {
				// Last bucket is overflow: estimate midpoint as 1.5x last bound
				midpoint = bounds[len(bounds)-1] * 1.5
			}
			totalSum += float64(count) * midpoint
		}

		sum := totalSum // avoid pointer-to-temp issue

		req.ResourceMetrics[0].ScopeMetrics[0].Metrics = append(req.ResourceMetrics[0].ScopeMetrics[0].Metrics,
			&otlpmetrics.Metric{
				Name:        "KloudMate.request.duration",
				Description: "Request duration distribution (cumulative histogram)",
				Unit:        "ms",
				Data: &otlpmetrics.Metric_Histogram{
					Histogram: &otlpmetrics.Histogram{
						DataPoints: []*otlpmetrics.HistogramDataPoint{
							{
								TimeUnixNano:      now,
								StartTimeUnixNano: timestamp10sAgo,
								Count:             totalCount,
								Sum:               &sum,
								ExplicitBounds:    bounds,
								BucketCounts:      bucketCounts,
								Attributes: []*commonpb.KeyValue{
									{
										Key:   "endpoint",
										Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "/api/v1/data"}},
									},
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
