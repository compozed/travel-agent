package models

import (
	"fmt"
	. "io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Name   string  `yaml:"name"`   // Supporting both JSON and YAML.
	Envs   []Env   `yaml:"envs"`   // Supporting both JSON and YAML.
	Groups []Group `yaml:"groups"` // Supporting both JSON and YAML.
}

type Env struct {
	Name      string   `yaml:"name"`       // Supporting both JSON and YAML.
	DependsOn []string `yaml:"depends_on"` // Supporting both JSON and YAML.
}

type Group struct {
	Name string `yaml:"name"` // Supporting both JSON and YAML.
}

func (f *Env) GetDependsOn() string {
	return fmt.Sprintf("[%s]", strings.Join(f.DependsOn, ","))
}

func (f *Env) GetDependsOnArray() []string {
	return f.DependsOn
}

func (e *Env) HasDependencies() bool {
	if e.DependsOn == nil || len(e.DependsOn) == 0 {
		return false
	} else {
		return true
	}
}

func Load(y []byte) (Config, error) {
	var config Config
	err := yaml.Unmarshal(y, &config)

	for index, env := range config.Envs {
		if len(env.DependsOn) == 0 {
			config.Envs[index].DependsOn = nil
		}
	}

	return config, err
}

func LoadFromFile(path string) (Config, error) {
	y, err := ReadFile(path)
	if err != nil {
		panic(err)
	}

	return Load(y)
}
