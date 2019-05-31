package main

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestMonitor_CallsNetFactoryWithRulesTimeout(t *testing.T) {
	expectedTimeout := 5 * time.Second
	rule := Rule{Website: EnvVariable("http://example.com"), Timeout: expectedTimeout}
	factory := func(timeout time.Duration) NetworkClient {
		assert.Equal(t, expectedTimeout, timeout, "monitor(_,_) expected to call the network factory with timeout %v; got %v", expectedTimeout, timeout)
		return &http.Client{}
	}
	monitor(rule, factory)
}

func TestMonitor_IncrementsErrorCounter_WhenRuleSpecifiesUnavailableCounter(t *testing.T) {
	reg := prometheus.NewRegistry()
	reg.MustRegister(MonitoringErrorCount)
	defer reg.Unregister(MonitoringErrorCount)

	rule := Rule{Website: EnvVariable("http://example.com"), Timeout: time.Second * 1, Counter: "Counter", Name: "Test"}
	factory := func(timeout time.Duration) NetworkClient {
		return &MockNetClient{statusCode: 300}
	}
	monitor(rule, factory)
	CheckPrometheusCounterVec(t, reg, MonitoringErrorCount, 1,
		ExpectationLabelPair{LabelName: checkerErrorCountLabel, LabelValue: invalidCounter},
	)
}

func TestMonitor_LogsStatusCode408_WhenCallTimesout(t *testing.T) {
	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "Namespace",
		Subsystem: "Subsystem",
		Name:      "Name",
		Help:      "Help",
	}, []string{"website", "status_code"})
	reg := prometheus.NewRegistry()
	reg.MustRegister(counter)
	defer reg.Unregister(counter)
	Counters["Counter"] = counter

	rule := Rule{Website: EnvVariable("http://example.com"), Timeout: time.Second * 1, Counter: "Counter", Name: "Test"}
	factory := func(timeout time.Duration) NetworkClient {
		return &MockNetClient{
			err: &MockTimeoutError{error: fmt.Errorf("internal timeout error")},
		}
	}
	monitor(rule, factory)
	CheckPrometheusCounterVec(t, reg, counter, 1,
		ExpectationLabelPair{LabelName: "website", LabelValue: rule.Website.String()},
		ExpectationLabelPair{LabelName: "status_code", LabelValue: "408"},
	)
}

func TestMonitor_LogsStatusCode500_WhenNetworkCallErrors(t *testing.T) {
	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "Namespace",
		Subsystem: "Subsystem",
		Name:      "Name",
		Help:      "Help",
	}, []string{"website", "status_code"})
	reg := prometheus.NewRegistry()
	reg.MustRegister(counter)
	defer reg.Unregister(counter)
	Counters["Counter"] = counter

	rule := Rule{Website: EnvVariable("http://example.com"), Timeout: time.Second * 1, Counter: "Counter", Name: "Test"}
	factory := func(timeout time.Duration) NetworkClient {
		return &MockNetClient{
			err: fmt.Errorf("network error"),
		}
	}
	monitor(rule, factory)
	CheckPrometheusCounterVec(t, reg, counter, 1,
		ExpectationLabelPair{LabelName: "website", LabelValue: rule.Website.String()},
		ExpectationLabelPair{LabelName: "status_code", LabelValue: "500"},
	)
}

func TestMonitor_LogsStatusCode_WhenResponseStatusCodeNotWithin2XX(t *testing.T) {
	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "Namespace",
		Subsystem: "Subsystem",
		Name:      "Name",
		Help:      "Help",
	}, []string{"website", "status_code"})
	reg := prometheus.NewRegistry()
	reg.MustRegister(counter)
	defer reg.Unregister(counter)
	Counters["Counter"] = counter

	rule := Rule{Website: EnvVariable("http://example.com"), Timeout: time.Second * 1, Counter: "Counter", Name: "Test"}
	factory := func(timeout time.Duration) NetworkClient {
		return &MockNetClient{statusCode: 300}
	}
	monitor(rule, factory)
	CheckPrometheusCounterVec(t, reg, counter, 1,
		ExpectationLabelPair{LabelName: "website", LabelValue: rule.Website.String()},
		ExpectationLabelPair{LabelName: "status_code", LabelValue: "300"},
	)
}

/*

 */

/********************************************************************************************************************************/
//////// Mocks
/********************************************************************************************************************************/
type MockTimeoutError struct {
	error
}

func (e MockTimeoutError) Timeout() bool {
	return true
}

func (e MockTimeoutError) Temporary() bool {
	return true
}

func (e MockTimeoutError) Error() string {
	return ""
}

type MockNetClient struct {
	err           error
	statusCode    int
	headCallCount int
	getCallCount  int
}

func (c *MockNetClient) Head(url string) (resp *http.Response, err error) {
	c.headCallCount += 1
	err = c.err
	if c.err != nil {
		return
	}

	resp.StatusCode = c.statusCode
	return
}

func (c *MockNetClient) Get(url string) (resp *http.Response, err error) {
	c.getCallCount += 1
	err = c.err
	if c.err != nil {
		return
	}

	resp = &http.Response{StatusCode: c.statusCode}
	return
}
