package main

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type ExpectationLabelPair struct {
	LabelName  string
	LabelValue string
}

// CheckPrometheusCounterVec is a helper method that checks that prometheus counter
// has the expected value for the expected label on the passed in registry
func CheckPrometheusCounterVec(t *testing.T, reg *prometheus.Registry, counter *prometheus.CounterVec, expectedValue float64, expectedLabels ...ExpectationLabelPair) {
	metricFamilies, err := reg.Gather()
	if err != nil {
		t.Fatalf("Unable to gather prometheus metrics: %+v", err)
	}

	if len(metricFamilies) < 1 ||
		len(metricFamilies[0].Metric) < 1 ||
		len(metricFamilies[0].Metric[0].Label) < 1 {
		metricCount := 0
		if len(metricFamilies) > 0 {
			metricCount = len(metricFamilies[0].Metric)
		}

		t.Fatalf("Unable to gather the metrics from prometheus.\n\tExpected 1 MetricFamilies; got: %d.\n\tExpected 1 Metric; got: %d", len(metricFamilies), metricCount)
	}

	var metricCounterValue float64 = 0.00
	metric := metricFamilies[0].Metric[0]

ExpectedLabels:
	for _, expectedLabelPair := range expectedLabels {
		log.Debugf("attempting to match expected label pair: %v", expectedLabelPair)
		for _, gotLabel := range metric.GetLabel() {
			log.Debugf("\tgot prometheus label: %s", gotLabel.GetName())
			if gotLabel.GetName() == expectedLabelPair.LabelName {
				metricCounterValue = metric.Counter.GetValue()
				gotValue := gotLabel.GetValue()
				if gotValue != expectedLabelPair.LabelValue {
					t.Fatalf(`Prometheus counter %+v expected [name: "%s" value:"%s"]; got [%s]`, counter, expectedLabelPair.LabelName, expectedLabelPair.LabelValue, gotLabel)
				}
				continue ExpectedLabels
			}
		}
	}

	switch {
	case metricCounterValue == 0.00:
		t.Fatalf("Prometheus counter %+v not found in the registry.", counter)
	case metricCounterValue != expectedValue:
		t.Fatalf("Prometheus counter %+v expected value %f; got %f", counter, expectedValue, metricCounterValue)
	}
}

// CheckPrometheusCounter is a helper method that checks that prometheus counter
// has the expected value on the passed in registry
func CheckPrometheusCounter(reg *prometheus.Registry, counter prometheus.Counter, expectedValue float64, t *testing.T) {
	metricFamilies, err := reg.Gather()
	if err != nil {
		t.Fatalf("Unable to gather prometheus metrics: %+v", err)
	}

	if len(metricFamilies) < 1 ||
		len(metricFamilies[0].Metric) < 1 {
		t.Fatalf("Unable to gather the metrics from prometheus.\n\tExpected 1 MetricFamilies; got: %d.\n\tExpected 1 Metric; got: %d", len(metricFamilies), len(metricFamilies[0].Metric))
	}

	metric := metricFamilies[0].Metric[0]
	if gotValue := metric.Counter.GetValue(); gotValue != expectedValue {
		t.Fatalf("Prometheus counter %+v expected count %f; got %f", counter, expectedValue, gotValue)
	}
}
