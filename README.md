Travel Agent
============

Works together with concourse to manage concourse manifest templates for multiple environment.

It allows to work with dependencies between.

## Installing

Make sure that your go environment is correctly set up in your workstation.

  # Pull dependencies
  go get github.com/onsi/ginkgo/ginkgo
  go get github.com/benbjohnson/ego/cmd/ego
  go get -d github.com/compozed/travel-agent

  # Make travel agent cli available (this asumes ~/bin is set in your $PATH)
  ln -s $GOPATH/src/compozed/travel-agent/bin/travel-agent ~/bin/.

## Usage

### Target

Its sets your concourse target

    ./travel-agent target CONCOURSE_IP:PORT

### bootstrap(TBD)

Generates travel agent structure, it also upgrades when newer versions of travel-agent get released

    ./travel-agent target CONCOURSE_IP:PORT

### Book

* Tests your **ci/manifest/manifest.ego**(Template that generates manifest) against **ci/manifest/assets/dev|prod**
* It generates manifest if tests passed
* If spruce_secrets.yml is provided it tries to spruce merge with generated manifest
* If generated manifest exists it tries to deploy to concourse


    ./trave-agent book travel-agent-config.yml spruce_secrets.yml

**travel-agent-config.yml** descrives environments that your pipeline supports:

    name: cf
    envs:
    - name: dev
    - name: test
      depends_on: dev
    - name: prod
      depends_on: test


**spruce_secrets.yml** configs that get merge after manifest generation
