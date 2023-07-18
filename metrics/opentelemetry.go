package metrics

import (
	"context"
	"fmt"
	"time"

	"github.com/goto/meteor/agent"
	"github.com/goto/meteor/config"
	"github.com/goto/meteor/recipe"
	"github.com/goto/salt/log"
	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/contrib/samplers/probability/consistent"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/encoding/gzip"
)

const gracePeriod = 5 * time.Second

// OtelMonitor represents the otel monitor.
type OtelMonitor struct {
	recipeDuration    metric.Int64Histogram
	extractorRetries  metric.Int64Counter
	assetsExtracted   metric.Int64Counter
	processorDuration metric.Int64Histogram
	sinkRetries       metric.Int64Counter
	sinkDuration      metric.Int64Histogram
}

func NewOTLP(ctx context.Context, cfg config.Config, logger *log.Logrus, appVersion string) (*OtelMonitor, func(), error) {
	if !cfg.OtelEnabled {
		return nil, noOp, nil
	}

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithOS(),
		resource.WithHost(),
		resource.WithProcess(),
		resource.WithContainer(),
		resource.WithProcessRuntimeName(),
		resource.WithProcessRuntimeVersion(),
		resource.WithAttributes(
			semconv.ServiceName(cfg.AppName),
			semconv.ServiceVersion(appVersion),
		))
	if err != nil {
		return nil, nil, fmt.Errorf("create resource: %w", err)
	}

	shutdownMetric, err := initGlobalMetrics(ctx, res, cfg, logger)
	if err != nil {
		return nil, nil, err
	}

	shutdownTracer, err := initGlobalTracer(ctx, res, cfg, logger)
	if err != nil {
		shutdownMetric()
		return nil, nil, err
	}

	shutdownProviders := func() {
		shutdownTracer()
		shutdownMetric()
	}

	if err := host.Start(); err != nil {
		shutdownProviders()
		return nil, nil, err
	}

	if err := runtime.Start(); err != nil {
		shutdownProviders()
		return nil, nil, err
	}

	// init meters
	meter := otel.Meter("")
	recipeDuration, err := meter.Int64Histogram("meteor.recipe.duration", metric.WithUnit("ms"))
	if err != nil {
		return nil, nil, err
	}

	processorDuration, err := meter.Int64Histogram("meteor.processor.duration", metric.WithUnit("ms"))
	if err != nil {
		return nil, nil, err
	}

	sinkDuration, err := meter.Int64Histogram("meteor.sink.duration", metric.WithUnit("ms"))
	if err != nil {
		return nil, nil, err
	}

	extractorRetries, err := meter.Int64Counter("meteor.extractor.retries")
	if err != nil {
		return nil, nil, err
	}

	assetsExtracted, err := meter.Int64Counter("meteor.assets.extracted")
	if err != nil {
		return nil, nil, err
	}

	sinkRetries, err := meter.Int64Counter("meteor.sink.retries")
	if err != nil {
		return nil, nil, err
	}

	return &OtelMonitor{
		recipeDuration:    recipeDuration,
		processorDuration: processorDuration,
		sinkDuration:      sinkDuration,
		extractorRetries:  extractorRetries,
		assetsExtracted:   assetsExtracted,
		sinkRetries:       sinkRetries,
	}, shutdownProviders, nil
}

func initGlobalMetrics(ctx context.Context, res *resource.Resource, cfg config.Config, logger *log.Logrus) (func(), error) {
	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(cfg.OtelCollectorAddr),
		otlpmetricgrpc.WithCompressor(gzip.Name),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("create metric exporter: %w", err)
	}

	reader := sdkmetric.NewPeriodicReader(
		exporter,
		sdkmetric.WithInterval(15*time.Second),
	)

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(reader),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(provider)

	return func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), gracePeriod)
		defer cancel()
		if err := provider.Shutdown(shutdownCtx); err != nil {
			logger.Error("otlp metric-provider failed to shutdown", "err", err)
		}
	}, nil
}

func initGlobalTracer(ctx context.Context, res *resource.Resource, cfg config.Config, logger *log.Logrus) (func(), error) {
	exporter, err := otlptrace.New(ctx, otlptracegrpc.NewClient(
		otlptracegrpc.WithEndpoint(cfg.OtelCollectorAddr),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithCompressor(gzip.Name),
	))
	if err != nil {
		return nil, fmt.Errorf("create trace exporter: %w", err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(consistent.ProbabilityBased(cfg.OtelTraceSampleProbability)),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exporter)),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), gracePeriod)
		defer cancel()
		if err := tracerProvider.Shutdown(shutdownCtx); err != nil {
			logger.Error("otlp trace-provider failed to shutdown", "err", err)
		}
	}, nil
}

func noOp() {}

func getSliceStringPluginNames(prs []recipe.PluginRecipe) []string {
	var res []string
	for _, pr := range prs {
		res = append(res, pr.Name)
	}

	return res
}

// RecordRun records a run behavior
func (m *OtelMonitor) RecordRun(ctx context.Context, run agent.Run) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("recipe_name", run.Recipe.Name))

	m.recipeDuration.Record(ctx,
		int64(run.DurationInMs),
		metric.WithAttributes(
			attribute.String("extractor", run.Recipe.Source.Name),
			attribute.StringSlice("processors", getSliceStringPluginNames(run.Recipe.Processors)),
			attribute.StringSlice("sinks", getSliceStringPluginNames(run.Recipe.Sinks)),
			attribute.Bool("success", run.Success),
		))

	m.extractorRetries.Add(ctx, int64(run.ExtractorRetries))

	m.assetsExtracted.Add(ctx,
		int64(run.AssetsExtracted),
		metric.WithAttributes(
			attribute.String("extractor", run.Recipe.Source.Name),
			attribute.StringSlice("processors", getSliceStringPluginNames(run.Recipe.Processors)),
			attribute.StringSlice("sinks", getSliceStringPluginNames(run.Recipe.Sinks)),
		))
}

// RecordPlugin records a individual plugin behavior in a run
func (m *OtelMonitor) RecordPlugin(ctx context.Context, pluginInfo agent.PluginInfo) {
	switch pluginInfo.PluginType {
	case "sink":
		m.sinkDuration.Record(ctx,
			time.Since(pluginInfo.StartTime).Milliseconds(),
			metric.WithAttributes(
				attribute.String("sink", pluginInfo.PluginName),
			))
		m.sinkRetries.Add(ctx,
			pluginInfo.Retries,
			metric.WithAttributes(
				attribute.String("sink", pluginInfo.PluginName),
				attribute.Int64("batch_size", int64(pluginInfo.BatchSize)),
			))

	case "processor":
		m.processorDuration.Record(ctx,
			time.Since(pluginInfo.StartTime).Milliseconds(),
			metric.WithAttributes(
				attribute.String("processor", pluginInfo.PluginName),
			))
	}
}
