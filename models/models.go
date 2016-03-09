package models

import (
	"fmt"
	"gopkg.in/yaml.v2"
	. "io/ioutil"
	"strings"
)

type Config struct {
	Name string `yaml:"name"` // Supporting both JSON and YAML.
	Envs []Env  `yaml:"envs"` // Supporting both JSON and YAML.
}

type Env struct {
	Name      string   `yaml:"name"`       // Supporting both JSON and YAML.
	DependsOn []string `yaml:"depends_on"` // Supporting both JSON and YAML.
}

func (f *Env) GetDependsOn() string {
	return fmt.Sprintf("[%s]", strings.Join(f.DependsOn, ","))
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
