package telemetry

import (
	"context"
	"errors"
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

type telemetryProvider struct {
	config        config.TelemetryConfig
	traceProvider *trace.TracerProvider
	meterProvider *metric.MeterProvider
}

type TelemetryProviderInterface interface {
	Shutdown(ctx context.Context) error
}

func New(config config.TelemetryConfig) (TelemetryProviderInterface, error) {
	res, err := newResource(config)
	if err != nil {
		return nil, err
	}

	traceProvider, err := newTraceProvider(context.Background(), config, res)
	if err != nil {
		return nil, err
	}

	meterProvider, err := newMeterProvider(context.Background(), config, res)
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(traceProvider)
	otel.SetMeterProvider(meterProvider)

	return &telemetryProvider{
		config:        config,
		traceProvider: traceProvider,
		meterProvider: meterProvider,
	}, nil
}

func newResource(config config.TelemetryConfig) (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(config.ServiceName),
			semconv.ServiceVersion(config.ServiceVersion),
		))
}

func newTraceProvider(ctx context.Context, config config.TelemetryConfig, res *resource.Resource) (*trace.TracerProvider, error) {
	options := []otlptracegrpc.Option{}
	if config.OtelExporterOtlpEndpoint != "" {
		options = append(options, otlptracegrpc.WithEndpoint(config.OtelExporterOtlpEndpoint))
	}

	if config.OtelExporterOtlpInsecure {
		options = append(options, otlptracegrpc.WithInsecure())
	}

	traceExporter, err := otlptracegrpc.New(ctx, options...)
	if err != nil {
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter, trace.WithBatchTimeout(time.Second)),
		trace.WithResource(res),
	)
	return traceProvider, nil
}

func newMeterProvider(ctx context.Context, config config.TelemetryConfig, res *resource.Resource) (*metric.MeterProvider, error) {
	options := []otlpmetricgrpc.Option{}
	if config.OtelExporterOtlpEndpoint != "" {
		options = append(options, otlpmetricgrpc.WithEndpoint(config.OtelExporterOtlpEndpoint))
	}

	if config.OtelExporterOtlpInsecure {
		options = append(options, otlpmetricgrpc.WithInsecure())
	}

	metricExp, err := otlpmetricgrpc.New(ctx, options...)
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(metricExp, metric.WithInterval(3*time.Second))),
	)
	return meterProvider, nil
}

func (o *telemetryProvider) Shutdown(ctx context.Context) error {
	var err error
	err = errors.Join(err, o.traceProvider.Shutdown(ctx))
	err = errors.Join(err, o.meterProvider.Shutdown(ctx))
	return err
}
