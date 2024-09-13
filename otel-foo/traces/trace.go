package traces

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// type myTracer struct {
// 	tracer trace.Tracer
// }

var MyTracer trace.Tracer

func InitMyTracer(serviceName string) {
	// return myTracer{
	// 	tracer: otel.Tracer(serviceName),
	// }
	MyTracer = otel.Tracer(serviceName)
}
