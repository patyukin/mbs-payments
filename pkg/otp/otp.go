package otp

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/prometheus"
	metricTlm "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	traceTlm "go.opentelemetry.io/otel/trace"
	"log"
)

type TelemetryWrapper struct {
	MeterProvider  *metric.MeterProvider
	TracerProvider *trace.TracerProvider
	Meter          metricTlm.Meter
	Tracer         traceTlm.Tracer
}

func newExporter() (sdktrace.SpanExporter, error) {
	exp, err := jaeger.New(
		jaeger.WithAgentEndpoint(
			jaeger.WithAgentHost("localhost"),
			jaeger.WithAgentPort("6831"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Jaeger exporter: %w", err)
	}

	return exp, nil
}

// NewTelemetryWrapper инициализирует MeterProvider и TracerProvider
func NewTelemetryWrapper() (*TelemetryWrapper, error) {
	promExporter, err := prometheus.New()
	if err != nil {
		log.Fatalf("Ошибка при создании экспортёра Prometheus: %v", err)
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(promExporter),
		metric.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				attribute.String("service.name", "my-service"), // Передаем KeyValue
			),
		),
	)

	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("ExampleService"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	exp, err := newExporter()
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)

	otel.SetMeterProvider(meterProvider)
	otel.SetTracerProvider(traceProvider)

	meter := meterProvider.Meter("my-service")
	tracer := otel.Tracer("my-service")

	return &TelemetryWrapper{
		MeterProvider:  meterProvider,
		TracerProvider: traceProvider,
		Tracer:         tracer,
		Meter:          meter,
	}, nil
}

// CreateCounter создает новый счётчик
func (tw *TelemetryWrapper) CreateCounter(name, description string) (metricTlm.Int64Counter, error) {
	counter, err := tw.Meter.Int64Counter(
		name,
		metricTlm.WithDescription(description),
	)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при создании счётчика: %v", err)
	}

	return counter, nil
}

// RecordCounter увеличивает значение счётчика
func (tw *TelemetryWrapper) RecordCounter(ctx context.Context, counter metricTlm.Int64Counter, value int64) {
	counter.Add(ctx, value)
}

// StartTrace создает новый span для трассировки
func (tw *TelemetryWrapper) StartTrace(ctx context.Context, name string) (context.Context, traceTlm.Span) {
	return tw.Tracer.Start(ctx, name)
}
