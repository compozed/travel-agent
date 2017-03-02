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
	var groups []Group

	BeforeEach(func() {
		envs = nil
		envs = append(envs, Env{Name: "dev"})
		envs = append(envs, Env{Name: "prod", DependsOn: []string{"dev"}})

		groups = nil
		groups = append(groups, Group{Name: "platform"})

		config.Envs = envs
		config.Name = "FOO"
		config.Groups = groups
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

	Describe("Env", func() {
		Describe("HasDependencies?", func() {

			It("should be true when there are dependencies", func() {
				var env = Env{Name: "prod"}
				env.DependsOn = append(env.DependsOn, "dev")

				Expect(env.HasDependencies()).Should(BeTrue())
			})

			It("should be false when there not any dependencies", func() {
				var env = Env{Name: "dev"}

				Expect(env.HasDependencies()).Should(BeFalse())
			})
		})

		Describe("GetDependsOn", func() {
			It("should return all dependencies as a string", func() {
				var env = Env{Name: "prod"}
				env.DependsOn = append(env.DependsOn, "dev")
				env.DependsOn = append(env.DependsOn, "test")

				Expect(env.GetDependsOn()).Should(Equal("[dev,test]"))
			})
		})

		Describe("GetDependsOnArray", func() {
			It("should return the array of dependencies", func() {
				var env = Env{Name: "prod"}
				env.DependsOn = append(env.DependsOn, "dev")
				env.DependsOn = append(env.DependsOn, "test")

				dependencies := []string{"dev", "test"}

				Expect(env.GetDependsOnArray()).Should(Equal(dependencies))
			})
		})
	})

	Describe("Group", func() {
		Describe("Get a groups", func() {
			It("should return a group name", func() {
				var err error
				var expected Config

				expected, err = LoadFromFile("example.yml")
				Expect(err).ShouldNot(HaveOccurred())

				Expect(expected.Groups[0].Name).Should(Equal("platform"))
			})
		})
	})
})
