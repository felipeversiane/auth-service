package telemetry

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/felipeversiane/auth-service/internal/infra/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

type telemetry struct {
	config        config.TelemetryConfig
	traceProvider *trace.TracerProvider
	meterProvider *metric.MeterProvider
}

type TelemetryInterface interface {
	Shutdown(ctx context.Context) error
}

func New(config config.TelemetryConfig) (TelemetryInterface, error) {
	slog.Info("initializing telemetry")

	res, err := newResource(config)
	if err != nil {
		slog.Error("failed to create resource", "error", err)
		return nil, err
	}

	traceProvider, err := newTraceProvider(context.Background(), config, res)
	if err != nil {
		slog.Error("failed to create trace provider", "error", err)
		return nil, err
	}

	meterProvider, err := newMeterProvider(context.Background(), config, res)
	if err != nil {
		slog.Error("failed to create meter provider", "error", err)
		return nil, err
	}

	otel.SetTracerProvider(traceProvider)
	otel.SetMeterProvider(meterProvider)

	slog.Info("telemetry initialized successfully")

	return &telemetry{
		config:        config,
		traceProvider: traceProvider,
		meterProvider: meterProvider,
	}, nil
}

func newResource(config config.TelemetryConfig) (*resource.Resource, error) {
	slog.Info("creating telemetry resource")

	res, err := resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(config.ServiceName),
			semconv.ServiceVersion(config.ServiceVersion),
		))
	if err != nil {
		slog.Error("failed to create telemetry resource", "error", err)
		return nil, err
	}

	return res, nil
}

func newTraceProvider(ctx context.Context, config config.TelemetryConfig, res *resource.Resource) (*trace.TracerProvider, error) {
	slog.Info("setting up trace provider", slog.String("otel_endpoint", config.OtelExporterOtlpEndpoint))

	options := []otlptracegrpc.Option{}
	if config.OtelExporterOtlpEndpoint != "" {
		options = append(options, otlptracegrpc.WithEndpoint(config.OtelExporterOtlpEndpoint))
	}

	if config.OtelExporterOtlpInsecure {
		options = append(options, otlptracegrpc.WithInsecure())
	}

	traceExporter, err := otlptracegrpc.New(ctx, options...)
	if err != nil {
		slog.Error("failed to initialize trace exporter", "error", err)
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter, trace.WithBatchTimeout(time.Second)),
		trace.WithResource(res),
	)

	slog.Info("trace provider initialized successfully")
	return traceProvider, nil
}

func newMeterProvider(ctx context.Context, config config.TelemetryConfig, res *resource.Resource) (*metric.MeterProvider, error) {
	slog.Info("setting up meter provider", slog.String("otel_endpoint", config.OtelExporterOtlpEndpoint))

	options := []otlpmetricgrpc.Option{}
	if config.OtelExporterOtlpEndpoint != "" {
		options = append(options, otlpmetricgrpc.WithEndpoint(config.OtelExporterOtlpEndpoint))
	}

	if config.OtelExporterOtlpInsecure {
		options = append(options, otlpmetricgrpc.WithInsecure())
	}

	metricExp, err := otlpmetricgrpc.New(ctx, options...)
	if err != nil {
		slog.Error("failed to initialize metric exporter", "error", err)
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(metricExp, metric.WithInterval(3*time.Second))),
	)

	slog.Info("meter provider initialized successfully")
	return meterProvider, nil
}

func (o *telemetry) Shutdown(ctx context.Context) error {
	slog.Info("shutting down telemetry services")

	var err error
	err = errors.Join(err, o.traceProvider.Shutdown(ctx))
	err = errors.Join(err, o.meterProvider.Shutdown(ctx))

	if err != nil {
		slog.Error("error during telemetry shutdown", "error", err)
	} else {
		slog.Info("telemetry shutdown completed successfully")
	}

	return err
}
