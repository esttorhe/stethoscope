package main

type ConfiguredCounter struct {
	Namespace string   `yaml:"namespace"`
	Subsystem string   `yaml:"subsystem"`
	Name      string   `yaml:"name"`
	Help      string   `yaml:"help"`
	Labels    []string `yaml:"labels"`
}
