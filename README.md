Travel Agent
============

Works together with concourse.ci to manage pipeline manifests.

## Goals:

**DRY pipeline manifests:** When pipelines perform the same task or deployment on multiple enviroments, manifests start getting large and repetitive.
Travel-agent addresses this issue by turning the pipeline manifest into a dynamic template

## Installing

Make sure that your go environment is correctly set up on your workstation.

    # Pull dependencies
    go get github.com/onsi/ginkgo/ginkgo
    go get github.com/onsi/gomega
    go get github.com/benbjohnson/ego/cmd/ego
    go get -d github.com/compozed/travel-agent/manifest

    # Make travel agent cmd available (this asumes ~/bin is set in your $PATH)
    ln -s $GOPATH/src/github.com/compozed/travel-agent/bin/travel-agent $GOPATH/bin/.


## Usage

### Target

Travel-Agent sets your concourse target

    ./travel-agent target CONCOURSE_IP:PORT

### bootstrap

Generates travel agent structure in `ci/manifest`

    ./travel-agent bootstrap PIPELINE_NAME 

### Book

    ./travel-agent book [TRAVEL_AGENT_CONFIG] [SPRUCE_SECRET_YAML]

1. Upgrades travel agent project when a newer versions is available locally 
2. Tests your **ci/manifest/manifest.ego** against **ci/manifest/assets/***
3. If all tests pass and `TRAVEL_AGENT_CONFIG` was provided, it will try to generate your manifest
4. If `SPRUCE_SECRET_YAML` is provided, it tries to spruce merge with the newly generated manifest
5. If a generated manifest exists, it tries to deploy to concourse

#### TRAVEL_AGENT_CONFIG

example:

    name: FOO
    git_project: https://github.com/compozed/travel-agent-example.git
    # git_project: /full/path/to/your/travel-agent-project
    envs:
    - name: dev
    - name: prod
      depends_on:
      - dev
        name: cf

