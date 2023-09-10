package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	defaultRegistry prometheus.Registerer
	defaultGatherer prometheus.Gatherer
)

func init() {
	defaultRegistry, defaultGatherer = newServiceRegistriy()
	MustRegister(DatabaseQueryRequestDuration)
}

// newServiceRegistriy returns new registry for service metrics.
// NOTE(nk2ge5k): automatically register collectors for runtime and process
// metrics, if any errors detected, while registering, panics.
func newServiceRegistriy() (prometheus.Registerer, prometheus.Gatherer) {
	gatherer := prometheus.NewRegistry()

	{ // register process metrics collector
		processCollector := collectors.NewProcessCollector(
			collectors.ProcessCollectorOpts{})
		gatherer.MustRegister(processCollector)
	}

	{ // register runtime metrics collector
		runtimeCollector := collectors.NewGoCollector()
		gatherer.MustRegister(runtimeCollector)
	}

	return gatherer, gatherer
}

// Register allows to register new metrics through global service registry
func Register(collector prometheus.Collector) error {
	return defaultRegistry.Register(collector)
}

// MustRegister allows to register new metrics through global service registry.
// Panics in case of error.
func MustRegister(collector prometheus.Collector) {
	if err := Register(collector); err != nil {
		panic(err)
	}
}

// Handler returns http handler for retrieving metrics from the server
func Handler() http.Handler {
	return promhttp.InstrumentMetricHandler(
		defaultRegistry, promhttp.HandlerFor(defaultGatherer, promhttp.HandlerOpts{}))
}
