package models_test

import (
	. "github.com/compozed/travel-agent/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

var _ = Describe("Models", func() {
	var envs []Env

	BeforeEach(func() {
		envs = nil
		envs = append(envs, Env{Name: "dev"})
		envs = append(envs, Env{Name: "prod", DependsOn: "a-devs-job"})
	})

	Describe("Load", func() {
		It("supports yaml confs", func() {
			var err error
			var y []byte
			var expected []Env

			y, err = yaml.Marshal(envs)
			Expect(err).ShouldNot(HaveOccurred())

			expected, err = Load(y)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(expected).Should(Equal(envs))
		})
	})

	Describe("LoadFromFile", func() {
		It("supports yaml confs", func() {
			var err error
			var expected []Env

			expected, err = LoadFromFile("example.yml")
			Expect(err).ShouldNot(HaveOccurred())

			Expect(expected).Should(Equal(envs))
		})
	})
})
