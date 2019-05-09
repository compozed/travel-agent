Travel Agent
============

Works together with [concourse](https://concourse-ci.org) to manage pipeline manifests. This tool can be useful when 
having 1 pipeline that deploys to multiple environments.

## Goals:

**DRY pipeline manifests:** When a pipelines perform the same steps to deploy in multiple environments, manifests start 
getting large and repetitive.  Travel-agent addresses this issue by turning the pipeline manifest into a dynamic template.

## Installing

    go get -d github.com/compozed/travel-agent/manifest
    ln -s $GOPATH/src/github.com/compozed/travel-agent/bin/travel-agent $GOPATH/bin/.

## Usage

### Target concourse

    ./travel-agent target CONCOURSE_IP:PORT

### Generate project

    cd PROJECT_NAME
    travel-agent init

[Customizing pipeline template](/docs/customizing_pipeline_templates.md)

### Configure pipeline


example `travel-agent.yml`:

    name: PIPELINE_NAME
    git_project: https://github.com/ORG/PROJECT_NAME.git
    # git_project: local/path/to/PROJECT_NAME
    envs:
    - name: dev
    - name: prod
      depends_on:
      - dev

[Configuring pipeline deployment](/docs/configuring_pipeline_deployment.md)

### Set concourse pipeline

    travel-agent book path/to/travel-agent.yml settingy1.yml

1. It will try to generate your manifest base on your `path/to/travel-agent.yml` config file
1. If `SPRUCE_SETTINGS_PATH` is provided, it tries to spruce merge with the newly generated manifest
1. It tries to deploy to concourse

## Contributing

See our [CONTRIBUTING](CONTRIBUTING.md) section for more information.


## License

The gem is available as open source under the terms of the [MIT License](http://opensource.org/licenses/MIT).
