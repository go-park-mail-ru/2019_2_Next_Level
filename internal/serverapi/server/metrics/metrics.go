package metrics

import "github.com/prometheus/client_golang/prometheus"

var Hits = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "hits",
}, []string{"status", "path"})
