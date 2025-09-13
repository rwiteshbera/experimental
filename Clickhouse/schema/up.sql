CREATE TABLE IF NOT EXISTS default.metrics_v4 ON CLUSTER chcluster
(
    workspaceId LowCardinality(String) CODEC(ZSTD(1)),
    fingerprint UInt64 CODEC(ZSTD(1)), 
    metric LowCardinality(String) CODEC(ZSTD(1)),
    description String CODEC(ZSTD(1)),
    unit LowCardinality(String) CODEC(ZSTD(1)),
    serviceName LowCardinality(String) CODEC(ZSTD(1)),
    timestamp_micro DateTime64(6) CODEC(Delta(8), ZSTD(1)),
    metric_type Enum8(
        'gauge' = 1,
        'sum' = 2,
        'histogram' = 3,
        'summary' = 4
    ) CODEC(ZSTD(1)),
    temporality Enum8(
        'unspecified' = 0,
        'cumulative' = 1,
        'delta' = 2
    ) CODEC(ZSTD(1)),
    is_monotonic UInt8 CODEC(ZSTD(1)),
    value Nullable(Float64) CODEC(Gorilla, ZSTD(1)),
    count Nullable(UInt64) CODEC(T64, ZSTD(1)),
    sum Nullable(Float64) CODEC(Gorilla, ZSTD(1)),
    buckets Nested (
        le Float64,                            
        count UInt64
    ) CODEC(ZSTD(1)),

    attributes Map(LowCardinality(String), String) CODEC(ZSTD(1)),
    exemplars Array(Tuple(
        spanId String,
        value Float64,
        traceId String,
        timestamp DateTime64(6),
        attributes Map(LowCardinality(String), String)
    )) CODEC(ZSTD(1)),
    _ttl Nullable(DateTime) DEFAULT NULL,
    INDEX idx_serviceName serviceName TYPE set(1000) GRANULARITY 1
)
ENGINE = MergeTree
PARTITION BY (toYYYYMM(timestamp_micro))
ORDER BY (workspaceId, metric, timestamp_micro)
SETTINGS index_granularity = 8192;



CREATE TABLE IF NOT EXISTS default.metrics_queue_v4 ON CLUSTER chcluster
(
    workspaceId String,
    fingerprint UInt64,
    metric String,
    description String,
    unit String,
    serviceName String,
    timestamp DateTime64(6),
    metric_type Enum8('gauge' = 1, 'sum' = 2, 'histogram' = 3, 'summary' = 4),
    temporality Enum8('unspecified' = 0, 'cumulative' = 1, 'delta' = 2),
    is_monotonic UInt8,
    value Float64,
    count UInt64,
    sum Float64,
    buckets Nested (le Float64, count UInt64),
    attributes Map(String, String),
    exemplars Array(Tuple(
        spanId String,
        value Float64,
        traceId String,
        timestamp DateTime64(6),
        attributes Map(String, String)
    )),
    _ttl DateTime,
)
ENGINE = Kafka('kafka:29092', 'metrics_v4', 'clickhouse-metrics_v4', 'JSONEachRow')
SETTINGS
    kafka_poll_timeout_ms = 500,
    kafka_num_consumers = 12,
    kafka_flush_interval_ms = 2000,
    kafka_max_block_size = 65536,
    kafka_thread_per_consumer = 1,
    kafka_poll_max_batch_size = 65536;


CREATE MATERIALIZED VIEW  IF NOT EXISTS default.mv_metrics_v4 ON CLUSTER chcluster
TO default.metrics_v4
AS
SELECT
    workspaceId,
    fingerprint,
    metric,
    description,
    unit,
    serviceName,
    timestamp as timestamp_micro,
    metric_type,
    temporality,
    is_monotonic,
    value,
    count,
    sum,
    arrayMap(s -> CAST(s AS Float64), buckets.le) AS `buckets.le`,
    buckets.count,
    attributes,
    exemplars,
    _ttl
FROM default.metrics_queue_v4;