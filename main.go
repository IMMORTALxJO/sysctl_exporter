package main

import (
	"flag"
	"net/http"

	sysctl "github.com/immortalxjo/go-sysctl"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	logLevel      = flag.String("log-level", "info", "Verbosity of logging")
	listenAddress = flag.String("listen-address", ":9141", "Address to listen on for telemetry")
	includeRegex  = flag.String("include", ".*", "RegExp for sysctl parameters")
	excludeRegex  = flag.String("exclude", "", "RegExp for skipping sysctl parameters")
	metricsPrefix = flag.String("metrics-prefix", "sysctl", "Prefix of prometheus metrics")
	procSysPath   = flag.String("proc-sys-path", "/proc/sys/", "Path to /proc/sys/ directory")
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	flag.Parse()
	switch *logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "warn", "warning":
		log.SetLevel(log.WarnLevel)
	}

	// Set /proc/sys/ path for go-syslog
	sysctlBase := *procSysPath
	if sysctlBase[len(sysctlBase)-1:] != "/" {
		sysctlBase = sysctlBase + "/"
	}
	sysctl.SysctlBase = sysctlBase

	exporter := &Exporter{
		includeRegex: *includeRegex,
		excludeRegex: *excludeRegex,
		prefix:       *metricsPrefix,
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
