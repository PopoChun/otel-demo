package metrics

import (
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var ApiPingRequestCount metric.Int64Counter
var ApiBarRequestCount metric.Int64Counter
var CommonLabels []attribute.KeyValue // display labels on prometheus metrics

type myMeter struct {
	meter metric.Meter
}

func InitMyMeter(serviceName string) myMeter {
	return myMeter{
		meter: otel.Meter(serviceName),
	}
}

func (my myMeter) Instrument() {
	CommonLabels = []attribute.KeyValue{
		attribute.String("company", "wit"),
		attribute.String("author", "someone"),
	}

	ApiPingRequestCount, _ = my.meter.Int64Counter(
		fmt.Sprintf("%s/%s", "fooooo_service_ping", "request_counts"),
		metric.WithDescription("The number of ping requests."),
	)

	ApiBarRequestCount, _ = my.meter.Int64Counter(
		fmt.Sprintf("%s/%s", "fooooo_service_bar", "request_counts"),
		metric.WithDescription("The number of bar requests."),
	)
}
