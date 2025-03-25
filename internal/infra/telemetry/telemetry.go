package telemetry

import (
	"context"
	"errors"
	"time"

	"github.com/felipeversiane/auth-service/internal/infra/config"
	"github.com/felipeversiane/auth-service/internal/infra/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.uber.org/zap"
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
	logger.Info("Initializing telemetry")

	res, err := newResource(config)
	if err != nil {
		logger.Error("Failed to create resource", zap.Error(err))
		return nil, err
	}

	traceProvider, err := newTraceProvider(context.Background(), config, res)
	if err != nil {
		logger.Error("Failed to create trace provider", zap.Error(err))
		return nil, err
	}

	meterProvider, err := newMeterProvider(context.Background(), config, res)
	if err != nil {
		logger.Error("Failed to create meter provider", zap.Error(err))
		return nil, err
	}

	otel.SetTracerProvider(traceProvider)
	otel.SetMeterProvider(meterProvider)

	logger.Info("Telemetry initialized successfully")

	return &telemetry{
		config:        config,
		traceProvider: traceProvider,
		meterProvider: meterProvider,
	}, nil
}

func newResource(config config.TelemetryConfig) (*resource.Resource, error) {
	logger.Info("Creating telemetry resource")

	res, err := resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(config.ServiceName),
			semconv.ServiceVersion(config.ServiceVersion),
		))
	if err != nil {
		logger.Error("Failed to create telemetry resource", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func newTraceProvider(ctx context.Context, config config.TelemetryConfig, res *resource.Resource) (*trace.TracerProvider, error) {
	logger.Info("Setting up trace provider", zap.String("otel_endpoint", config.OtelExporterOtlpEndpoint))

	options := []otlptracegrpc.Option{}
	if config.OtelExporterOtlpEndpoint != "" {
		options = append(options, otlptracegrpc.WithEndpoint(config.OtelExporterOtlpEndpoint))
	}

	if config.OtelExporterOtlpInsecure {
		options = append(options, otlptracegrpc.WithInsecure())
	}

	traceExporter, err := otlptracegrpc.New(ctx, options...)
	if err != nil {
		logger.Error("Failed to initialize trace exporter", zap.Error(err))
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter, trace.WithBatchTimeout(time.Second)),
		trace.WithResource(res),
	)

	logger.Info("Trace provider initialized successfully")
	return traceProvider, nil
}

func newMeterProvider(ctx context.Context, config config.TelemetryConfig, res *resource.Resource) (*metric.MeterProvider, error) {
	logger.Info("Setting up meter provider", zap.String("otel_endpoint", config.OtelExporterOtlpEndpoint))

	options := []otlpmetricgrpc.Option{}
	if config.OtelExporterOtlpEndpoint != "" {
		options = append(options, otlpmetricgrpc.WithEndpoint(config.OtelExporterOtlpEndpoint))
	}

	if config.OtelExporterOtlpInsecure {
		options = append(options, otlpmetricgrpc.WithInsecure())
	}

	metricExp, err := otlpmetricgrpc.New(ctx, options...)
	if err != nil {
		logger.Error("Failed to initialize metric exporter", zap.Error(err))
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(metricExp, metric.WithInterval(3*time.Second))),
	)

	logger.Info("Meter provider initialized successfully")
	return meterProvider, nil
}

func (o *telemetry) Shutdown(ctx context.Context) error {
	logger.Info("Shutting down telemetry services")

	var err error
	err = errors.Join(err, o.traceProvider.Shutdown(ctx))
	err = errors.Join(err, o.meterProvider.Shutdown(ctx))

	if err != nil {
		logger.Error("Error during telemetry shutdown", zap.Error(err))
	} else {
		logger.Info("Telemetry shutdown completed successfully")
	}

	return err
}
