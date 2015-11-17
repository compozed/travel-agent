package templates_test

import (
	. "github.allstate.com/CompoZedPlatform/travel-agent/ci"
	. "github.allstate.com/CompoZedPlatform/travel-agent/models"

	"bytes"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var _ = Describe("Nfs", func() {
	envs := []Env{
		Env{Name: "development"},
		Env{Name: "staging"},
	}

	Context("when having multiple environments", func() {
		envsLocal := envs
		for index, env := range envs {
			envLocal := env
			indexLocal := index

			BeforeEach(func() {
				if indexLocal > 0 {
					envLocal.DependsOn = fmt.Sprintf("nfs-%s", envsLocal[indexLocal-1].Name)
				}
			})

			It(fmt.Sprintf("applies template dependencies correctly for %s", envLocal.Name), func() {
				var err error
				var buf bytes.Buffer
				err = JobTmpl(&buf, envLocal)
				if err != nil {
				}
				actualJob := make(map[interface{}]interface{})
				err = yaml.Unmarshal(buf.Bytes(), &actualJob)
				Expect(actualJob).ShouldNot(BeNil())

				expectedJob := ReadYAML(fmt.Sprintf("assets/nfs_%s.yml", envsLocal[indexLocal].Name))
				Expect(actualJob).Should(Equal(expectedJob))
			})

			It(fmt.Sprintf("should render nfs-%s", envLocal.Name), func() {
				var err error
				var buf bytes.Buffer

				err = JobTmpl(&buf, envLocal)
				if err != nil {
				}

				actualJob := make(map[interface{}]interface{})
				err = yaml.Unmarshal(buf.Bytes(), &actualJob)

				Expect(actualJob).ShouldNot(BeNil())
			})
		}
	})
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
