package models

import (
	"gopkg.in/yaml.v2"
	. "io/ioutil"
)

type Config struct {
	Name string `yaml:"name"` // Supporting both JSON and YAML.
	Envs []Env  `yaml:"envs"` // Supporting both JSON and YAML.
}

type Env struct {
	Name      string `yaml:"name"`       // Supporting both JSON and YAML.
	DependsOn string `yaml:"depends_on"` // Supporting both JSON and YAML.
}

func Load(y []byte) (Config, error) {
	var config Config
	err := yaml.Unmarshal(y, &config)
	return config, err
}

func LoadFromFile(path string) (Config, error) {
	y, err := ReadFile(path)
	if err != nil {

		panic(err)
	}

	return Load(y)

}
