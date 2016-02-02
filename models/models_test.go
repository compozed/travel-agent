package models_test

import (
	. "github.com/compozed/travel-agent/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

var _ = Describe("Models", func() {
	var config Config
	var envs []Env

	BeforeEach(func() {
		envs = nil
		envs = append(envs, Env{Name: "dev"})
		envs = append(envs, Env{Name: "prod", DependsOn: "dev"})

		config.Envs = envs
		config.Name = "FOO"
	})

	Describe("Load", func() {
		It("supports yaml confs", func() {
			var err error
			var y []byte
			var expected Config

			y, err = yaml.Marshal(config)
			Expect(err).ShouldNot(HaveOccurred())

			expected, err = Load(y)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(expected).Should(Equal(config))
		})
	})

	Describe("LoadFromFile", func() {
		It("supports yaml confs", func() {
			var err error
			var expected Config

			expected, err = LoadFromFile("example.yml")
			Expect(err).ShouldNot(HaveOccurred())

			Expect(expected).Should(Equal(config))
		})
	})
})
