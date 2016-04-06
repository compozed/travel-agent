Travel Agent
============

Works together with concourse.ci to manage pipeline manifests.

## Goals:

**DRY pipeline manifests:** When pipelines needs to perform the exact task or deployment on multiple enviroments pipeline's manifests start gettting bigger.
This is address by travel-agent by turning the pipeline manifest into a `ego`


## Installing

Make sure that your go environment is correctly set up in your workstation.

    # Pull dependencies
    go get github.com/onsi/ginkgo/ginkgo
    go get github.com/onsi/gomega
    go get github.com/benbjohnson/ego/cmd/ego
    go get -d github.com/compozed/travel-agent/manifest

    # Make travel agent cli available (this asumes ~/bin is set in your $PATH)
    ln -s $GOPATH/src/github.com/compozed/travel-agent/bin/travel-agent ~/bin/.


## Usage

### Target

Travel-Agent sets your concourse target

    ./travel-agent target CONCOURSE_IP:PORT

### bootstrap

Generates travel agent structure in `ci/manifest`

    ./travel-agent bootstrap PIPELINE_NAME 

### Book

    ./travel-agent book [TRAVEL_AGENT_CONFIG] [SPRUCE_SECRET_YAML]

1. Upgrades travel agent project when a newer versions is available in your system
2. Tests your **ci/manifest/manifest.ego** against **ci/manifest/assets/***
3. If tests passed and `TRAVEL_AGENT_CONFIG` was provided it will try to generate your manifest
4. If `SPRUCE_SECRET_YAML` is provided, it tries to spruce merge with the generated manifest
5. If generated manifest exists, it tries to deploy to concourse

#### TRAVEL_AGENT_CONFIG

example:

    name: cf
    envs:
    - name: dev
    - name: test
      depends_on: dev
    - name: prod
      depends_on: test


### Writing concourse manifest templates

Travel agent relays on `ego` to generate pipeline manifests. After a bootstrap you will have a `ci/manifest/manifest.ego` that you can use as an starting point.
