package templates_test

import (
	. "github.allstate.com/CompoZedPlatform/concourse-workspace/models"
	. "github.allstate.com/CompoZedPlatform/concourse-workspace/templates"

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
		Env{Name: "production"},
	}

	Context("when having multiple environments", func() {
		envsLocal := envs
		for index, env := range envs {
			envLocal := env
			indexLocal := index

			BeforeEach(func() {
				if indexLocal > 0 {
					envsLocal[indexLocal].DependsOn = fmt.Sprintf("nfs-%s", envsLocal[indexLocal-1].Name)
				}
			})

			It(fmt.Sprintf("applies template dependencies correctly for %s", envLocal.Name), func() {
				var err error
				var buf bytes.Buffer
				err = NfsTmpl(&buf, envsLocal)
				if err != nil {
				}
				actualJob := GetJob(fmt.Sprintf("nfs-%s", envsLocal[indexLocal].Name), buf)
				Expect(actualJob).ShouldNot(BeNil())

				expectedJob := ReadYAML(fmt.Sprintf("assets/nfs_%s.yml", envsLocal[indexLocal].Name))
				Expect(actualJob).Should(Equal(expectedJob))
			})

			It(fmt.Sprintf("should render nfs-%s", envLocal.Name), func() {
				var err error
				var buf bytes.Buffer
				err = NfsTmpl(&buf, envsLocal)
				if err != nil {
				}
				actualJob := GetJob(fmt.Sprintf("nfs-%s", envLocal.Name), buf)
				Expect(actualJob).ShouldNot(BeNil())

				// Put step for deployment has the environment name in the name
				deploymentName := actualJob["plan"].([]interface{})[3].(map[interface{}]interface{})["put"]
				Expect(deploymentName).Should(Equal(fmt.Sprintf("nfs-%s-deployment", envLocal.Name)))

				// spruce merge merges the correct stub environment
				spruceFiles := actualJob["plan"].([]interface{})[2].(map[interface{}]interface{})["config"].(map[interface{}]interface{})["params"].(map[interface{}]interface{})["files"]
				Expect(spruceFiles).Should(Equal(fmt.Sprintf("templates/nfs/manifest.yml deployments/%s/nfs/stub.yml", envLocal.Name)))
			})
		}
	})
})

func GetJob(jobName string, manifest bytes.Buffer) map[interface{}]interface{} {
	m := make(map[interface{}]interface{})
	err := yaml.Unmarshal(manifest.Bytes(), &m)
	if err != nil {
	}

	jobs := m["jobs"].([]interface{})
	var job map[interface{}]interface{}
	for i := 0; job == nil && i < len(jobs); i++ {
		job = jobs[i].(map[interface{}]interface{})
		if job["name"] != jobName {
			job = nil
		}
	}
	return job
}
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
