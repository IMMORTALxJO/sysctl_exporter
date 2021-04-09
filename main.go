package main

import (
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var (
	logLevel      = flag.String("log-level", "info", "Verbosity of logging")
	listenAddress = flag.String("listen-address", ":9141", "Address to listen on for telemetry")
	grepPattern   = flag.String("pattern", ".*", "Regexp for sysctl parameters")
	skipPattern   = flag.String("skip-pattern", "", "Regexp for skipping sysctl parameters")
	metricsPrefix = flag.String("metrics-prefix", "sysctl", "Prefix of prometheus metrics")
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	flag.Parse()
	switch *logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "warn", "warning":
		log.SetLevel(log.WarnLevel)
	}
}

func main() {
	exporter := &Exporter{
		grepPattern: *grepPattern,
		skipPattern: *skipPattern,
		prefix:      *metricsPrefix,
	}
	prometheus.MustRegister(exporter)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Sysctl Exporter</title></head>
             <body>
             <h1>Sysctl Exporter</h1>
             <p><a href='/metrics'>Metrics</a></p>
             </body>
             </html>`))
	})
	log.Fatal(http.ListenAndServe(*listenAddress, nil))

}
