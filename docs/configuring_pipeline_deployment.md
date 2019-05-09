# Configuring pipeline deployment

#### Global features

When setting a pipeline with travel agent you can toggle certain jobs or resources 
depending on a flag on your configuration. This can be done through global features.

`travel-agent.yml`:

    ---
    name: PIPELINE_NAME
    features:
      slack-notifications: true 
      ...

#### Environment features

when setting a pipeline with travel agent you can toggle certain steps for a env 
job depending on a flag on your configuration. This features will only apply 
to some environments.

`travel-agent.yml`:

    ---
    envs:
    ...
    - name: prod
      depends_on:
      - dev
      features:
        backup: true 
