package main

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
)

func main() {
	// Configure logging level
	if err := log.Base().SetLevel("DEBUG"); err != nil {
		panic(err.Error())
	}
	configuration, err := LoadConfiguration()
	if err != nil {
		log.Errorf("Unable to load configuration. %s", err.Error())
		return
	}

	// Expose metrics to prometheus
	http.Handle("/metrics", promhttp.Handler())

	// Load the rules
	for _, rule := range configuration.Rules {
		go startMonitoring(rule)
	}

	log.Fatal(http.ListenAndServe(":7000", nil))
}

func startMonitoring(rule Rule) {
	log.Debugf("starting monitoring rule: %+v", rule)
	go monitor(rule, NewNetClient)
	for range time.Tick(time.Duration(rule.Interval)) {
		monitor(rule, NewNetClient)
	}
}
