package main

import (
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/common/log"
)

type NetClientFactory func(time.Duration) NetworkClient
type NetworkClient interface {
	Head(string) (*http.Response, error)
	Get(string) (*http.Response, error)
}

func monitor(rule Rule, clientFactory NetClientFactory) {
	timeout := time.Duration(rule.Timeout)
	netClient := clientFactory(timeout)

	// Use HEAD by default unles rule specifies otherwise
	var method func(string) (*http.Response, error) = netClient.Head
	if !rule.UseHead {
		log.Debug("using GET")
		method = netClient.Get
	}

	res, err := method(rule.Website.String())
	switch {
	// An error occurred pinging the website
	case err != nil:
		handleError(err, rule)
	// Website replied but response is an error
	case res != nil && res.StatusCode < 200 || res.StatusCode > 299:
		logErrorCode(res.StatusCode, rule)
	// Website is up
	default:
		log.Infof("website '%s' is up. Checking again in %d minutes", rule.Website, rule.Interval.Round(time.Minute))
	}
}

func handleError(err error, rule Rule) {
	log.Errorf("Cannot load %s. %s", rule.Website, err.Error())
	statusCode := http.StatusInternalServerError
	if err, ok := err.(net.Error); ok && err.Timeout() {
		statusCode = http.StatusRequestTimeout
	}
	logErrorCode(statusCode, rule)
}

func logErrorCode(statusCode int, rule Rule) {
	if counter, ok := Counters[rule.Counter]; ok {
		log.Infof("\tupping counter %s for rule %s", rule.Counter, rule.Name)
		counter.WithLabelValues(rule.Website.String(), strconv.Itoa(statusCode)).Inc()
	} else {
		log.Errorf("\tinvalid counter specified «%s» for rule «%s»", rule.Counter, rule.Name)
		MonitoringErrorCount.WithLabelValues(invalidCounter).Inc()
	}
}
