package main

import (
	"io/ioutil"

	"github.com/prometheus/common/log"
	"gopkg.in/yaml.v2"
)

type MonitoringConfiguration struct {
	Rules []Rule `yaml:"rules"`
}

func LoadConfiguration() (config MonitoringConfiguration, err error) {
	// Parse the rules file
	fileName := "rules.yml"
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Errorf("unable to read %s file; %s", fileName, err.Error())
		return
	}

	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		log.Errorf("unable to parse %s file; %s", fileName, err.Error())
	}
	log.Debugf("loaded rules configuration: %+v", config)

	return
}
