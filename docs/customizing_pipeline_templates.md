# Ego Templates

Travel-Agent relays on [ego](https://github.com/benbjohnson/ego) as its template engine to build the [pipeline](https://concourse-ci.org/pipelines.html) configuration for concourse.

See pipeline example [here](/manifest/manifest.ego)
## Tags

Ego uses tags to embed golang code.

There are 2 tags available:

1. Execution tag `<% %>`
1. Execution and output tag `<%= %>`

### Execution only tag

Code gets executed but the result will not be included in the pipeline configuration file.

```
<% if { %>
<% } %>
```

### Execution and output tag

Code gets executed and the output gets included in the pipeline configuration file.

```
<%= env.Name %> 
```

## Golang and Travel-agent

### Conditionals 

```
<% if CONDITION { %>
<% } else { %>
<% } %>
```

### Iterators

```
<% for _, VAR_NAME := range ARRAY { %>
<% } %>
```

### Travel-agent helpers

Travel-agent allows you to retrive data from the travel-agent.yml config file through helpers.

```
<%! func ManifestTmpl(w io.Writer, config Config) error %>
<%% import . "github.com/compozed/travel-agent/models" %%>
```

### config.Envs

Retuns list of all the available environments.


### env.HasFeature()

Returns features available for an specific environment

```
<% if env.HasFeature("check_for_existence") { %>
```

### env.Feature()

```
'<%= env.Feature("grab_value") %>'
```

Pulls data for an specific environment.

```
name: <PIPELINE_NAME>
envs:
- name: <ENV_NAME>
  features:
    grab_value: foo
```
