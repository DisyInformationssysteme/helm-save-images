# Helm Save-Images Plugin

This is a Helm plugin to download an save Docker images for a given Helm Chart and an optional list of Helm Values files.

## Install

```shell
helm plugin install https://github.com/DisyInformationssysteme/helm-save-images
```

## Usage

```bash
Usage:
  helm save-images [chart] [flags]

Flags:
      --dry-run          print list of images only
  -h, --help             help for helm
  -f, --values strings   specify values in a YAML file or a URL (can specify multiple)
      --version string   chart version
```

## Build

Clone the repository into your `$GOPATH` and then build it.

```bash
git clone https://github.com/DisyInformationssysteme/helm-save-images.git
cd helm-save-images
go build
```

The above will install this plugin into your `$HELM_HOME/plugins` directory.

### Prerequisites

- You need to have [Go](http://golang.org) installed. Make sure to set `$GOPATH`

### Running Tests

Automated tests are implemented with [*testing*](https://golang.org/pkg/testing/).

To run all tests:

```bash
go test -v ./...
```

## Release

Set `GITHUB_TOKEN` and run:

```bash
goreleaser release
```
