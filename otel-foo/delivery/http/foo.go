package http

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/PopoChun/otel-demo/otel-foo/config"
	"github.com/PopoChun/otel-demo/otel-foo/metrics"
	"github.com/PopoChun/otel-demo/otel-foo/traces"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"
)

type fooHandler struct {
}

func NewFooHandler() *fooHandler {
	return &fooHandler{}
}

func (f *fooHandler) SayHello(c *gin.Context) {
	barServerApi := fmt.Sprintf("%s/%s", config.GetBarServerConfig().Host, "hello")
	req, err := http.NewRequestWithContext(c, "GET", barServerApi, nil)
	if err != nil {
		fmt.Printf("%s: %s", "failed to create http request", err)
	}
	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("%s: %s", "failed to do the request", err)
	}
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)
	c.JSON(http.StatusOK, gin.H{"message": string(data)})
}

func (f *fooHandler) Ping(c *gin.Context) {
	metrics.ApiPingRequestCount.Add(c.Request.Context(), 1, metrics.CommonLabels...)
	span := trace.SpanFromContext(c.Request.Context())
	span.SetAttributes(attribute.String("hanlder", "fooHandler"))
	span.AddEvent("test event calling /ping", trace.WithAttributes(attribute.Int("intVal", 1234), attribute.String("stringVal", "xoxo")))
	c.JSON(http.StatusOK, gin.H{"message": "Pong"})
}

func (f *fooHandler) GetBar(c *gin.Context) {
	metrics.ApiBarRequestCount.Add(c.Request.Context(), 1)

	member1, _ := baggage.NewMember("hello", "world")
	member2, _ := baggage.NewMember("asdf", "qwer")
	bag, _ := baggage.New(member1, member2)

	defaultCtx := baggage.ContextWithBaggage(context.Background(), bag)
	ctx, span := traces.MyTracer.Start(defaultCtx, "Request_To_Bar")

	// Trace an HTTP client by wrapping the transport
	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	// Make sure we pass the context to the request to avoid broken traces
	barServerApi := fmt.Sprintf("%s/%s", config.GetBarServerConfig().Host, "callbar")
	req, err := http.NewRequestWithContext(ctx, "GET", barServerApi, nil)
	if err != nil {
		fmt.Printf("%s: %s", "failed to create http request", err)
	}

	// All requests made with this client will create spans.
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("%s: %s", "failed to do the request", err)
	}
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)
	span.End()
	c.JSON(http.StatusOK, gin.H{"message": string(data)})
}
