package metrics

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
)

var ApiPingRequestCount instrument.Int64Counter
var ApiBarRequestCount instrument.Int64Counter
var CommonLabels []attribute.KeyValue // display labels on prometheus metrics

type myMeter struct {
	meter metric.Meter
}

func InitMyMeter(serviceName string) myMeter {
	return myMeter{
		meter: global.Meter(serviceName),
	}
}

func (my myMeter) Instrument() {
	CommonLabels = []attribute.KeyValue{
		attribute.String("company", "wit"),
		attribute.String("author", "john.chang"),
	}

	ApiPingRequestCount, _ = my.meter.Int64Counter(
		fmt.Sprintf("%s/%s", "fooooo_service_ping", "request_counts"),
		instrument.WithDescription("The number of ping requests."),
	)

	ApiBarRequestCount, _ = my.meter.Int64Counter(
		fmt.Sprintf("%s/%s", "fooooo_service_bar", "request_counts"),
		instrument.WithDescription("The number of bar requests."),
	)
}
