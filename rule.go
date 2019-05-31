package main

import (
	"fmt"
	"os"
	"time"

	"github.com/prometheus/common/log"
)

// EnvVariable is a wrapper around 'string' type.
// Provides custom 'YAML' unmarshaling. It tries to laod
// the environment variable specified in the 'yaml' file.
type EnvVariable string

// String returns a string representation of the `EnvVariable`
func (e EnvVariable) String() string {
	return string(e)
}

func (e *EnvVariable) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var asString string
	if err := unmarshal(&asString); err != nil {
		return err
	}

	if len(asString) < 1 {
		return fmt.Errorf("env variable cannot be empty")
	}

	envVar := os.Getenv(asString)
	log.Debugf("loaded env var «%s» -> %s", asString, envVar)
	*e = EnvVariable(envVar)
	return nil
}

// Rule represents how often (`Interval`) the system will ping `Website`.
// The ping will timeout after the specified `Timeout` upping the `Counter` every time
// the system is unable to reach the `Website`.
// The system will attempt to use `HEAD` method by default unles `UseHead` specifies otherwise
type Rule struct {
	Name     string        `yaml:"name"`
	Website  EnvVariable   `yaml:"website"`
	Timeout  time.Duration `yaml:"timeout"`
	Interval time.Duration `yaml:"interval"`
	Counter  string        `yaml:"counter"`
	UseHead  bool          `yaml:"use_head"`
}
