package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

type Config struct {
	serviceName    string
	environment    string
	jaegerEndpoint string
}

func NewConfig(serviceName, environment, jaegerEndpoint string) Config {
	return Config{
		serviceName:    serviceName,
		environment:    environment,
		jaegerEndpoint: jaegerEndpoint,
	}
}

func InitTracer(cfg Config) (func(context.Context) error, error) {
	// exporter
	traceExporter, err := newExporter(cfg.jaegerEndpoint)
	if err != nil {
		return nil, err
	}

	// tracer provider
	tracerProvider, err := newTracerProvider(cfg, traceExporter)
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(tracerProvider)

	// propagator
	propagator := newPropagator()
	otel.SetTextMapPropagator(propagator)

	return tracerProvider.Shutdown, nil
}

func newExporter(endpoint string) (sdktrace.SpanExporter, error) {
	return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)))
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTracerProvider(cfg Config, traceExporter sdktrace.SpanExporter) (*sdktrace.TracerProvider, error) {
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.serviceName),
			semconv.DeploymentEnvironmentKey.String(cfg.environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create the resource: %w", err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)

	return tracerProvider, nil
}
