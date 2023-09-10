package metrics

import "github.com/prometheus/client_golang/prometheus"

const namespace = "late"

var DatabaseQueryRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: "storage",
		Name:      "query_request_duration_seconds",
		Help:      "The amount of time, in seconds, it took to execute a database query.",
	},
	[]string{"query", "kind"},
)
