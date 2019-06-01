package main

import (
	"io/ioutil"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"gopkg.in/yaml.v2"
)

const (
	PrometheusCheckerNamespace  = "stethoscope"
	PrometheusCheckertSubsystem = "webhook"
	checkerErrorCountLabel      = "reason"

	invalidCounter = "invalid_counter"
)

// MetricsConfiguration holds the list of configured counters (read from `counters.yml`)
type MetricsConfiguration struct {
	Counters []ConfiguredCounter `yaml:"counters"`
}

var (
	// MonitoringErrorCount is the `prometheus counter` used when there's an internal
	// unexpected error while running the service.
	MonitoringErrorCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: PrometheusCheckerNamespace,
		Subsystem: PrometheusCheckertSubsystem,
		Name:      "error",
		Help:      "Number of errors while monitoring external websites",
	}, []string{checkerErrorCountLabel})

	// Counters maps counter names to actual `prometheus.CounterVec` instances
	Counters map[string]*prometheus.CounterVec = make(map[string]*prometheus.CounterVec)
)

func init() {
	prometheus.MustRegister(
		MonitoringErrorCount,
	)

	data, err := ioutil.ReadFile("counters.yml")
	if err != nil {
		log.Errorf("unable to read counters.yml file; %s", err.Error())
		return
	}

	configuration := MetricsConfiguration{}
	err = yaml.Unmarshal([]byte(data), &configuration)
	if err != nil {
		log.Errorf("unable to parse counters.yml file; %s", err.Error())
		return
	}

	log.Infof("parsed config\n%+v", configuration)

	for _, configCounter := range configuration.Counters {
		counter := prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: configCounter.Namespace,
			Subsystem: configCounter.Subsystem,
			Name:      configCounter.Name,
			Help:      configCounter.Help,
		}, configCounter.Labels)
		prometheus.MustRegister(counter)
		Counters[configCounter.Name] = counter
	}

	log.Infof("loaded counters\n%+v", Counters)
}
