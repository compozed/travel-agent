Travel Agent
============

Works together with **concourse.ci** to manage pipeline manifests.

## Goals:

**DRY pipeline manifests:** When pipelines perform the same task or deployment on multiple enviroments, manifests start getting large and repetitive.
Travel-agent addresses this issue by turning the pipeline manifest into a dynamic template

## Installing

Make sure that your go environment is correctly set up on your workstation.

    go get -d github.com/compozed/travel-agent/manifest
    ln -s $GOPATH/src/github.com/compozed/travel-agent/bin/travel-agent $GOPATH/bin/.

## Usage

### Target

Travel-Agent sets your concourse target

    ./travel-agent target CONCOURSE_IP:PORT

### bootstrap

Generates travel agent structure in `ci/manifest`

    cd YOUR_PROJECT
    ./travel-agent init

### Book

    ./travel-agent book [TRAVEL_AGENT_CONFIG_PATH] [SPRUCE_SETTINGS_PATH]

1. It will try to generate your manifest base on your `TRAVEL_AGENT_CONFIG`
1. If `SPRUCE_SETTINGS_PATH` is provided, it tries to spruce merge with the newly generated manifest
1. It tries to deploy to concourse

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


## Contributing

See our [CONTRUBUTING](CONTRIBUTING.md) section for more information.


## License

The gem is available as open source under the terms of the [MIT License](http://opensource.org/licenses/MIT).
