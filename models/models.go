package models

import (
	"fmt"
	. "io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Name      string     `yaml:"name"` // Supporting both JSON and YAML.
	Envs      []Env      `yaml:"envs"`
	Groups    []Group    `yaml:"groups"`
	Resources []Resource `yaml:"resources"`
}

type Env struct {
	Name      string                      `yaml:"name"`
	DependsOn []string                    `yaml:"depends_on"`
	Features  map[interface{}]interface{} `yaml:"features"`
}

type Resource struct {
	Name string `yaml:"name"`
}

type Group struct {
	Name string `yaml:"name"`
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

func (e *Env) HasFeature(feature string) bool {
	_, ok := e.Features[feature]
	return ok
}

func (e *Env) Feature(feature string) string {
	if v, ok := e.Features[feature]; ok {
		if v == nil {
			return ""
		}
		if _, ok := v.(map[interface{}]interface{}); ok {
			panic(fmt.Sprintf("Using a map as the value of a feature is currently unsupported (on feature '%s')\n", feature))
		}
		if _, ok := v.([]interface{}); ok {
			panic(fmt.Sprintf("Using an array as the value of feature requires using the FeaturesArray() function (on feature '%s')\n", feature))
		}
		return fmt.Sprintf("%v", v)
	} else {
		return ""
	}
}

func (e *Env) FeatureList(feature string) []string {
	if v, ok := e.Features[feature]; ok {
		if v == nil {
			return []string{}
		}
		features := []string{}
		if list, ok := v.([]interface{}); ok {
			for _, element := range list {
				if _, ok := element.(map[interface{}]interface{}); ok {
					panic(fmt.Sprintf("Using a map as the value of a feature list is currently unsupported (on feature '%s')\n", feature))
				}
				if _, ok := element.([]interface{}); ok {
					panic(fmt.Sprintf("Using an array as the value of feature list is currently unsupported (on feature '%s')\n", feature))
				}
				features = append(features, fmt.Sprintf("%v", element))
			}
			return features
		} else {
			panic(fmt.Sprintf("Tried to call FeatureList on '%s', but its value was not a list of things. Got '%#v' instead", feature, v))
		}
	} else {
		return []string{}
	}
}

func Load(y []byte) (Config, error) {
	var config Config
	err := yaml.Unmarshal(y, &config)

	for index, env := range config.Envs {
		if len(env.DependsOn) == 0 {
			config.Envs[index].DependsOn = nil
		}
		if env.Features == nil {
			config.Envs[index].Features = map[interface{}]interface{}{}
		}
	}
	if config.Resources == nil {
		config.Resources = []Resource{}
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
