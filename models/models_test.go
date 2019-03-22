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
	var resources []Resource

	BeforeEach(func() {
		envs = nil
		envs = append(envs, Env{Name: "dev", Features: map[interface{}]interface{}{}})
		envs = append(envs, Env{Name: "prod", DependsOn: []string{"dev"}, Features: map[interface{}]interface{}{}})

		groups = nil
		groups = append(groups, Group{Name: "platform"})

		resources = nil
		resources = append(resources, Resource{Name: "jumpbox"})

		config.Envs = envs
		config.Name = "FOO"
		config.Groups = groups
		config.Resources = resources
		config.Features = map[interface{}]interface{}{}
	})

	Describe("Load", func() {
		It("supports yaml confs", func() {
			var err error
			var y []byte
			var result Config

			y, err = yaml.Marshal(config)
			Expect(err).ShouldNot(HaveOccurred())

			result, err = Load(y)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(result).Should(Equal(config))
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
		It("should instantiate features if omitted from yaml", func() {
			var err error
			var expected Config
			expected, err = LoadFromFile("example.yml")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(expected.Envs[0].Features).ShouldNot(BeNil())
			Expect(expected.Envs[1].Features).ShouldNot(BeNil())
		})
	})

	Describe("Config", func() {
		Describe("HasFeature", func() {
			var config = Config{Features: map[interface{}]interface{}{"iaas": "aws"}}
			It("should return true if the env has a feature", func() {
				Expect(config.HasFeature("iaas")).Should(BeTrue())
			})
			It("should return false if the feature  is not present", func() {
				Expect(config.HasFeature("nonexistent")).Should(BeFalse())
			})
			It("should return false (and not crash) if there are no features at all", func() {
				var config = Config{}
				Expect(config.HasFeature("testFeature")).Should(Equal(false))
			})
		})
		Describe("Feature", func() {
			It("should return a stringified version of the feature", func() {
				var config = Env{Features: map[interface{}]interface{}{"testFeature": "string"}}
				Expect(config.Feature("testFeature")).Should(Equal("string"))
			})
			It("should stringify booleans", func() {
				var config = Config{Features: map[interface{}]interface{}{"testFeature": true}}
				Expect(config.Feature("testFeature")).Should(Equal("true"))
			})
			It("should stringify floats", func() {
				var config = Config{Features: map[interface{}]interface{}{"testFeature": 2.4}}
				Expect(config.Feature("testFeature")).Should(Equal("2.4"))
			})
			It("should stringify ints", func() {
				var config = Config{Features: map[interface{}]interface{}{"testFeature": 42}}
				Expect(config.Feature("testFeature")).Should(Equal("42"))
			})
			It("does not support maps as feature values", func() {
				testMaps := func() {
					var config = Config{Features: map[interface{}]interface{}{"testFeature": map[interface{}]interface{}{"stuff": "value"}}}
					config.Feature("testFeature")
				}
				Expect(testMaps).Should(Panic())
			})
			It("does not support arrays as feature values", func() {
				testArrays := func() {
					var config = Config{Features: map[interface{}]interface{}{"testFeature": []interface{}{"stuff", "more stuff"}}}
					config.Feature("testFeature")
				}
				Expect(testArrays).Should(Panic())
			})
			It("should stringify nulls as an empty string", func() {
				var config = Config{Features: map[interface{}]interface{}{"testFeature": nil}}
				Expect(config.Feature("testFeature")).Should(Equal(""))
			})
			It("should return an empty string if the feature isn't defined", func() {
				var config = Config{Features: map[interface{}]interface{}{}}
				Expect(config.Feature("testFeature")).Should(Equal(""))
			})
			It("should return an empty string (and not crash) if there are no features at all", func() {
				var config = Config{}
				Expect(config.Feature("testFeature")).Should(Equal(""))
			})
		})
		Describe("FeatureList", func() {
			It("should return a list stringified objects of the feature", func() {
				var config = Config{Features: map[interface{}]interface{}{"testFeature": []interface{}{"string"}}}
				Expect(config.FeatureList("testFeature")).Should(Equal([]string{"string"}))
			})
			It("should stringify booleans", func() {
				var config = Config{Features: map[interface{}]interface{}{"testFeature": []interface{}{true}}}
				Expect(config.FeatureList("testFeature")).Should(Equal([]string{"true"}))
			})
			It("should stringify floats", func() {
				var config = Config{Features: map[interface{}]interface{}{"testFeature": []interface{}{2.4}}}
				Expect(config.FeatureList("testFeature")).Should(Equal([]string{"2.4"}))
			})
			It("should stringify ints", func() {
				var config = Config{Features: map[interface{}]interface{}{"testFeature": []interface{}{42}}}
				Expect(config.FeatureList("testFeature")).Should(Equal([]string{"42"}))
			})
			It("does not configsupport maps as FeatureList values", func() {
				testMaps := func() {
					var config = Config{Features: map[interface{}]interface{}{"testFeature": map[interface{}]interface{}{"stuff": "value"}}}
					config.FeatureList("testFeature")
				}
				Expect(testMaps).Should(Panic())
			})
			It("does not support arrays as FeatureList element values", func() {
				testArrays := func() {
					var config = Config{Features: map[interface{}]interface{}{
						"testFeature": []interface{}{
							[]interface{}{"stuff", "more stuff"},
						}}}
					config.FeatureList("testFeature")
				}
				Expect(testArrays).Should(Panic())
			})
			It("should stringify nulls as an empty list", func() {
				var config = Config{Features: map[interface{}]interface{}{"testFeature": nil}}
				Expect(config.FeatureList("testFeature")).Should(Equal([]string{}))
			})
			It("should return an empty list if the feature isn't defined", func() {
				var config = Config{Features: map[interface{}]interface{}{}}
				Expect(config.FeatureList("testFeature")).Should(Equal([]string{}))
			})
			It("Should panic if the value found with FeatureList is not actually a list", func() {
				testValueNotList := func() {
					var config = Config{Features: map[interface{}]interface{}{"testFeature": true}}
					config.FeatureList("testFeature")
				}
				Expect(testValueNotList).Should(Panic())
			})
			It("should return an empty list (and not crash) if there are no features at all", func() {
				var config = Config{}
				Expect(config.FeatureList("testFeature")).Should(Equal([]string{}))
			})
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
		Describe("HasFeature", func() {
			var env = Env{Features: map[interface{}]interface{}{"iaas": "aws"}}
			It("should return true if the env has a feature", func() {
				Expect(env.HasFeature("iaas")).Should(BeTrue())
			})
			It("should return false if the feature  is not present", func() {
				Expect(env.HasFeature("nonexistent")).Should(BeFalse())
			})
			It("should return false (and not crash) if there are no features at all", func() {
				var env = Env{}
				Expect(env.HasFeature("testFeature")).Should(Equal(false))
			})
		})
		Describe("Feature", func() {
			It("should return a stringified version of the feature", func() {
				var env = Env{Features: map[interface{}]interface{}{"testFeature": "string"}}
				Expect(env.Feature("testFeature")).Should(Equal("string"))
			})
			It("should stringify booleans", func() {
				var env = Env{Features: map[interface{}]interface{}{"testFeature": true}}
				Expect(env.Feature("testFeature")).Should(Equal("true"))
			})
			It("should stringify floats", func() {
				var env = Env{Features: map[interface{}]interface{}{"testFeature": 2.4}}
				Expect(env.Feature("testFeature")).Should(Equal("2.4"))
			})
			It("should stringify ints", func() {
				var env = Env{Features: map[interface{}]interface{}{"testFeature": 42}}
				Expect(env.Feature("testFeature")).Should(Equal("42"))
			})
			It("does not support maps as feature values", func() {
				testMaps := func() {
					var env = Env{Features: map[interface{}]interface{}{"testFeature": map[interface{}]interface{}{"stuff": "value"}}}
					env.Feature("testFeature")
				}
				Expect(testMaps).Should(Panic())
			})
			It("does not support arrays as feature values", func() {
				testArrays := func() {
					var env = Env{Features: map[interface{}]interface{}{"testFeature": []interface{}{"stuff", "more stuff"}}}
					env.Feature("testFeature")
				}
				Expect(testArrays).Should(Panic())
			})
			It("should stringify nulls as an empty string", func() {
				var env = Env{Features: map[interface{}]interface{}{"testFeature": nil}}
				Expect(env.Feature("testFeature")).Should(Equal(""))
			})
			It("should return an empty string if the feature isn't defined", func() {
				var env = Env{Features: map[interface{}]interface{}{}}
				Expect(env.Feature("testFeature")).Should(Equal(""))
			})
			It("should return an empty string (and not crash) if there are no features at all", func() {
				var env = Env{}
				Expect(env.Feature("testFeature")).Should(Equal(""))
			})
		})
		Describe("FeatureList", func() {
			It("should return a list stringified objects of the feature", func() {
				var env = Env{Features: map[interface{}]interface{}{"testFeature": []interface{}{"string"}}}
				Expect(env.FeatureList("testFeature")).Should(Equal([]string{"string"}))
			})
			It("should stringify booleans", func() {
				var env = Env{Features: map[interface{}]interface{}{"testFeature": []interface{}{true}}}
				Expect(env.FeatureList("testFeature")).Should(Equal([]string{"true"}))
			})
			It("should stringify floats", func() {
				var env = Env{Features: map[interface{}]interface{}{"testFeature": []interface{}{2.4}}}
				Expect(env.FeatureList("testFeature")).Should(Equal([]string{"2.4"}))
			})
			It("should stringify ints", func() {
				var env = Env{Features: map[interface{}]interface{}{"testFeature": []interface{}{42}}}
				Expect(env.FeatureList("testFeature")).Should(Equal([]string{"42"}))
			})
			It("does not support maps as FeatureList values", func() {
				testMaps := func() {
					var env = Env{Features: map[interface{}]interface{}{"testFeature": map[interface{}]interface{}{"stuff": "value"}}}
					env.FeatureList("testFeature")
				}
				Expect(testMaps).Should(Panic())
			})
			It("does not support arrays as FeatureList element values", func() {
				testArrays := func() {
					var env = Env{Features: map[interface{}]interface{}{
						"testFeature": []interface{}{
							[]interface{}{"stuff", "more stuff"},
						}}}
					env.FeatureList("testFeature")
				}
				Expect(testArrays).Should(Panic())
			})
			It("should stringify nulls as an empty list", func() {
				var env = Env{Features: map[interface{}]interface{}{"testFeature": nil}}
				Expect(env.FeatureList("testFeature")).Should(Equal([]string{}))
			})
			It("should return an empty list if the feature isn't defined", func() {
				var env = Env{Features: map[interface{}]interface{}{}}
				Expect(env.FeatureList("testFeature")).Should(Equal([]string{}))
			})
			It("Should panic if the value found with FeatureList is not actually a list", func() {
				testValueNotList := func() {
					var env = Env{Features: map[interface{}]interface{}{"testFeature": true}}
					env.FeatureList("testFeature")
				}
				Expect(testValueNotList).Should(Panic())
			})
			It("should return an empty list (and not crash) if there are no features at all", func() {
				var env = Env{}
				Expect(env.FeatureList("testFeature")).Should(Equal([]string{}))
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

	Describe("Resource", func() {
		Describe("Get a resource", func() {
			It("should return a resourcename", func() {
				var err error
				var expected Config

				expected, err = LoadFromFile("example.yml")
				Expect(err).ShouldNot(HaveOccurred())

				Expect(expected.Resources[0].Name).Should(Equal("jumpbox"))
			})
		})
	})
})
