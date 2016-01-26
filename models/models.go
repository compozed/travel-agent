package models

import (
	"gopkg.in/yaml.v2"
	. "io/ioutil"
)

type Env struct {
	Name      string `yaml:"name"`       // Supporting both JSON and YAML.
	DependsOn string `yaml:"depends_on"` // Supporting both JSON and YAML.
}

func Load(y []byte) ([]Env, error) {

	var envs []Env
	err := yaml.Unmarshal(y, &envs)
	return envs, err
}

func LoadFromFile(path string) ([]Env, error) {
	y, err := ReadFile(path)
	if err != nil {

		panic(err)
	}

	return Load(y)

}
