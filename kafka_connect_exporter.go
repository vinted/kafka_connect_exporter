package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/vinted/kafka-connect-exporter/internal/app/kafka-connect-exporter/collector"
	"net/http"
	"os"
)

const (
	nameSpace  = "kafka_connect"
	version    = "dev"
	versionUrl = "https://github.com/vinted/kafka-connect-exporter"
)

var (
	showVersion   = flag.Bool("version", false, "show version and exit")
	listenAddress = flag.String("listen-address", ":8080", "Address on which to expose metrics.")
	metricsPath   = flag.String("telemetry-path", "/metrics", "Path under which to expose metrics.")
	scrapeURI     = flag.String("scrape-uri", "http://127.0.0.1:8080", "URI on which to scrape kafka connect.")
	user          = flag.String("user", "", "Optional username for authenticating to kafka-connect")
	pass          = flag.String("pass", "", "Optional password for authenticating to kafka-connect")
)

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Printf("kafka_connect_exporter\n url: %s\n version: %s\n", versionUrl, version)
		os.Exit(2)
	}

	log.Infoln("Starting kafka_connect_exporter")

	prometheus.Unregister(prometheus.NewGoCollector())
	prometheus.Unregister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	prometheus.MustRegister(collector.NewCollector(*scrapeURI, nameSpace, *user, *pass))

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, *metricsPath, http.StatusMovedPermanently)
	})

	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
