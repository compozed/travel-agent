package main_test

import (
	. "../manifest"
	"bytes"
	"fmt"
	. "github.com/compozed/travel-agent/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var _ = Describe("Manifest generation", func() {
	config, _ := LoadFromFile("assets/travel-agent.yml")

	for index, env := range config.Envs {
		envLocal := env
		indexLocal := index
		var actualManifest map[interface{}]interface{}

		JustBeforeEach(func() {
			var err error
			var buf bytes.Buffer

			if indexLocal > 0 {
				envLocal.DependsOn = config.Envs[indexLocal-1].Name
			}

			config := Config{config.Name, []Env{envLocal}}
			err = ManifestTmpl(&buf, config)
			Ω(err).ShouldNot(HaveOccurred())

			err = yaml.Unmarshal(buf.Bytes(), &actualManifest)
			if err != nil {
				fmt.Println("Generated Manifest:")
				fmt.Println(string(buf.Bytes()))
			}
			Ω(err).ShouldNot(HaveOccurred())

			Expect(actualManifest).ShouldNot(BeNil())
		})

		Describe("When rendering jobs", func() {
			expectedJobs := ReadYAML(fmt.Sprintf("assets/%s.yml", envLocal.Name))["jobs"].([]interface{})

			for _, job := range expectedJobs {
				localExpectedJob := job.(map[interface{}]interface{})

				It(fmt.Sprintf("Should render %s", localExpectedJob["name"]), func() {
					jobName := localExpectedJob["name"].(string)

					actualJob := GetJob(actualManifest, jobName)

					Expect(actualJob).ShouldNot(BeNil(), fmt.Sprintf("Manifest did not generate job \"%s\"", jobName))
					Expect(actualJob).Should(Equal(localExpectedJob))
				})
			}

		})

		Describe("When rendering resources", func() {
			expectedYaml := ReadYAML(fmt.Sprintf("assets/%s.yml", envLocal.Name))

			if expectedYaml["resources"] != nil {
				expectedResources := expectedYaml["resources"].([]interface{})

				for _, resource := range expectedResources {
					localExpectedResource := resource.(map[interface{}]interface{})

					It(fmt.Sprintf("Should render %s", localExpectedResource["name"]), func() {
						resourceName := localExpectedResource["name"].(string)

						actualResource := GetResource(actualManifest, resourceName)

						Expect(actualResource).Should(Equal(localExpectedResource))
					})
				}
			}
		})
	}
})

func ReadYAML(filepath string) map[interface{}]interface{} {
	res := make(map[interface{}]interface{})
	source, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(source, &res)
	if err != nil {
		panic(err)
	}
	return res
}

func GetJob(manifest map[interface{}]interface{}, jobName string) map[interface{}]interface{} {
	return GetItem(manifest, "jobs", jobName)
}

func GetResource(manifest map[interface{}]interface{}, resourceName string) map[interface{}]interface{} {
	return GetItem(manifest, "resources", resourceName)
}

func GetItem(manifest map[interface{}]interface{}, itemType string, itemName string) map[interface{}]interface{} {
	items := manifest[itemType].([]interface{})
	var item map[interface{}]interface{}
	for i := 0; item == nil && i < len(items); i++ {
		item = items[i].(map[interface{}]interface{})
		if item["name"] != itemName {
			item = nil
		}
	}
	return item
}
