package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"google.golang.org/grpc"

	"github.com/PopoChun/otel-demo/otel-foo/config"
	_httpDelivery "github.com/PopoChun/otel-demo/otel-foo/delivery/http"
	"github.com/PopoChun/otel-demo/otel-foo/metrics"
	"github.com/PopoChun/otel-demo/otel-foo/traces"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const serviceName = "foo-service"

func initProvider() func() {
	ctx := context.Background()
	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			// the service name used to display traces in backends
			// display {job="foo-service"} label on prometheus metrics
			semconv.ServiceNameKey.String(serviceName),
		))
	handleErr(err, "fail to create resource")

	otelColAddr, ok := os.LookupEnv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if !ok {
		otelColAddr = config.GetOtelCollectorConfig().Host
	}

	metricExp, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(otelColAddr))
	handleErr(err, "failed to create the collector metric exporter")

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(
				metricExp,
				sdkmetric.WithInterval(2*time.Second),
			),
		),
	)
	global.SetMeterProvider(meterProvider)

	traceClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(otelColAddr),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
	)
	// given 3 seconds to connect to traceClient
	sctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	traceExp, err := otlptrace.New(sctx, traceClient)
	handleErr(err, "Failed to create the collector trace exporter")

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(traceExp),
	)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tracerProvider)

	return func() {
		cxt, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		if err := traceExp.Shutdown(cxt); err != nil {
			otel.Handle(err)
		}
		if err := meterProvider.Shutdown(cxt); err != nil {
			otel.Handle(err)
		}
	}
}

func handleErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}

func main() {
	config.InitConf()

	shutdown := initProvider()
	defer shutdown()

	traces.InitMyTracer(fmt.Sprintf("%s-%s", serviceName, "tracer"))
	// meter := global.Meter(fmt.Sprintf("%s-%s", serviceName, "meter"))
	myMeter := metrics.InitMyMeter(fmt.Sprintf("%s-%s", serviceName, "meter"))
	myMeter.Instrument()

	r := gin.Default()
	r.Use(otelgin.Middleware(serviceName))

	fooHandler := _httpDelivery.NewFooHandler()
	r.GET("/ping", fooHandler.Ping)
	r.GET("/bar", fooHandler.GetBar)
	r.GET("/omg", fooHandler.SayHello)
	r.Run(config.GetHttpConfig().Port)
}
