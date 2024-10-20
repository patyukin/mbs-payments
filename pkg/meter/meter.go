package meter

import (
	"context"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

type MetricsWrapper struct {
	meterProvider metric.MeterProvider
}

func NewMetricsWrapper() *MetricsWrapper {
	config := prometheus.Config{}
	exporter, err := prometheus.New(config)
	if err != nil {
		log.Fatalf("Ошибка при создании экспортёра Prometheus: %v", err)
	}

	// Создаем MeterProvider на основе экспортёра
	meterProvider := metric.NewMeterProvider(metric.WithReader(exporter))
	otel.SetMeterProvider(meterProvider)

	// Получаем Meter для регистрации метрик
	meter := meterProvider.Meter("my-service")

	// Возвращаем структуру
	return &MetricsWrapper{
		meterProvider: meterProvider,
		meter:         meter,
	}
}

// StartPrometheusServer запускает сервер Prometheus для экспорта метрик
func (mw *MetricsWrapper) StartPrometheusServer() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Метрики доступны на http://localhost:2112/metrics")
		log.Fatal(http.ListenAndServe(":2112", nil))
	}()
}

// CreateCounter создает новый счётчик
func (mw *MetricsWrapper) CreateCounter(name, description string) metric.Int64Counter {
	counter, err := mw.meter.Int64Counter(
		name,
		metric.WithDescription(description),
	)
	if err != nil {
		log.Fatalf("Ошибка при создании счётчика: %v", err)
	}
	return counter
}

// CreateHistogram создает гистограмму
func (mw *MetricsWrapper) CreateHistogram(name, description string) metric.Float64Histogram {
	histogram, err := mw.meter.Float64Histogram(
		name,
		metric.WithDescription(description),
	)
	if err != nil {
		log.Fatalf("Ошибка при создании гистограммы: %v", err)
	}
	return histogram
}

// RecordCounter увеличивает значение счетчика
func (mw *MetricsWrapper) RecordCounter(ctx context.Context, counter metric.Int64Counter, value int64) {
	counter.Add(ctx, value)
}

// RecordHistogram записывает значение в гистограмму
func (mw *MetricsWrapper) RecordHistogram(ctx context.Context, histogram metric.Float64Histogram, value float64) {
	histogram.Record(ctx, value)
}
