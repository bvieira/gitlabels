# gitlabels

Create, update and/or remove labels from one or more repositories of an organization or user

## Getting Started

### Installing

Download the lastest version on github release

### Config

`gitlabels` is configured using a `YAML` file folowing the format below:

| attribute | type | mandatory | description|
| --------- |:----:| :---: | :----------|
| owner | string | :warning: | github owner |
| org | string | :warning: | github organization, has priority over `owner` |
| project-regex | string | :heavy_check_mark: | `regex` for project names |
| labels | map<string, label-config> | | collection of labels that will be added or updated on repositories |
| {label-config}.color | string | | label color, use hexadecimal colors
| {label-config}.description | string | | label description |
| rename | map<string, string> | | collection of labels name that will be renamed, it will use `labels` configs. `key` is the current name, `value` is the new name.
| delete | []string | | collection of labels name that will be removed

:warning: `owner` or `org` should exists!

#### Config sample

```yaml
owner: some-owner
project-regex: .*
labels:
  'bug':
    color: 801515
    description: this is bug
  'bugfix':
    description: this is bug
  'wontfix':
    color: 801515
rename:
  'rename-label': 'bug'
delete:
  - 'delete-label1'
  - 'delete-label2'
```

## Running

run `gitlabels-cli -h` to get the list of available parameters

```bash
gitlabels-cli -h
```

after configuring your `cfg.yaml`, run `gitlabels-cli`

```bash
#if cfg is not set, cfg.yaml will be used
gitlabels-cli -token github-token

#use -cfg to use a different config
gitlabels-cli -cfg cfg.sample.yaml -token github-token
```

## Development

### Makefile

run `make` to get the list of available actions

```bash
make
```

#### Make configs

| Variable | description|
| --------- | ----------|
| BUILDOS | build OS |
| BUILDARCH | build arch |
| ECHOFLAGS | flags used on echo |
| BUILDENVS | var envs used on build |
| BUILDFLAGS | flags used on build |

| Parameters | description|
| --------- | ----------|
| args | parameters that will be used on run |

```bash
#variables
BUILDOS="linux" BUILDARCH="amd64" make build

#parameters
make run args="--cfg cfg.sample.yaml"
```

### Build

```bash
make build
```

the binary will be created on `bin/$BUILDOS_$BUILDARCH/gitlabels-cli`

### Tests

```bash
make test
```

### Run

```bash
#without args
make run

#with args
make run args="--cfg cfg.sample.yaml"
```
