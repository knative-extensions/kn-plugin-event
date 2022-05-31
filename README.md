# kn-plugin-event

[![Mage](https://github.com/knative-sandbox/kn-plugin-event/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/knative-sandbox/kn-plugin-event/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/knative-sandbox/kn-plugin-event)](https://goreportcard.com/report/knative-sandbox/kn-plugin-event)
[![Releases](https://img.shields.io/github/release-pre/knative-sandbox/kn-plugin-event.svg?sort=semver)](https://github.com/knative-sandbox/kn-plugin-event/releases)
[![LICENSE](https://img.shields.io/github/license/knative-sandbox/kn-plugin-event.svg)](https://github.com/knative-sandbox/kn-plugin-event/blob/main/LICENSE)
[![Maturity level](https://img.shields.io/badge/Maturity%20level-ALPHA-red)](https://github.com/knative/community/tree/main/mechanics/MATURITY-LEVELS.md)

`kn-plugin-event` is a plugin of Knative Client, for managing cloud events from
command line.

## Description

With this plugin, you can build and send the cloud events to publicly available
addresses (URLs) or to Addressable (Kubernetes service, Knative service, broker,
channels, etc) inside of Kubernetes cluster.

## Usage

### Build an event

You could use `build` subcommand to create an event with a help of convenient
builder-like interface. It may be helpful to send that even later or use it in
other context.

#### Simplest event

```sh
$ kn event build -o yaml
data: {}
datacontenttype: application/json
id: 81a402a2-9c29-4c27-b8ed-246a253c9e58
source: kn-event/v0.4.0
specversion: "1.0"
time: "2021-10-15T10:42:57.713226203Z"
type: dev.knative.cli.plugin.event.generic
```

#### More complex example

```sh
$ kn event build \
    --field operation.type=local-wire-transfer \
    --field operation.amount=2345.40 \
    --field operation.from=87656231 \
    --field operation.to=2344121 \
    --field automated=true \
    --field signature='FGzCPLvYWdEgsdpb3qXkaVp7Da0=' \
    --type org.example.bank.bar \
    --id $(head -c 10 < /dev/urandom | base64 -w 0) \
    --output json
{
  "specversion": "1.0",
  "id": "RjtL8UH66X+UJg==",
  "source": "kn-event/v0.4.0",
  "type": "org.example.bank.bar",
  "datacontenttype": "application/json",
  "time": "2021-10-15T10:43:23.113187943Z",
  "data": {
    "automated": true,
    "operation": {
      "amount": "2345.40",
      "from": 87656231,
      "to": 2344121,
      "type": "local-wire-transfer"
    },
    "signature": "FGzCPLvYWdEgsdpb3qXkaVp7Da0="
  }
}
```

### Send an event

To send and event you should utilize `send` subcommand. The `send` command
leverages the same builder-like interface as `build` command does. You could
send an event to public address of your application by using `--to-url` option,
or to supported in-cluster resources by using `--to` option.

#### Sending to public address

To send event to public address you use `--to-url` option:

```sh
$ kn event send \
    --field player.id=6354aa60-ddb1-452e-8c13-24893667de20 \
    --field player.game=2345 \
    --field points=456 \
    --type org.example.gaming.foo \
    --to-url http://ce-api.foo.example.org/
```

> **NOTE**: All argument, except `--to-url` are optional. Use
> `kn event send --help` to see full usage information.

#### Sending to in-cluster resources

Sending event to in-cluster resources is done with a companion Job that is
deployed on your cluster. In that way `kn-event` and send your event even to
resources that are not publicly accessible.

Send an event to a Knative service `event-display` in namespace `default`:

```sh
$ kn event send \
    --type org.example.kn.ping \
    --id $(uuidgen) \
    --field event.type=test \
    --field event.data=98765 \
    --to Service:serving.knative.dev/v1:event-display
```

> **NOTE**: All argument, except `--to-url` are optional. Use
> `kn event send --help` to see full usage information.

## Install

You can download a pre-built version of `kn-plugin-event` from
[our release page](https://github.com/knative-sandbox/kn-plugin-event/releases)
. Choose the one that fits your platform.

When the download is ready, you should be ready to use `kn-plugin-event` as a
standalone binary. Check the available commands with:

```sh
kn-event-<OS>-<ARCH> --help
```

### Install to work with `kn` CLI

If you'd like to use the plugin with `kn` CLI, install the plugin by simply
copying the executable file under `kn` plugins directory as:

```sh
$ mkdir -p ~/.config/kn/plugins
$ cp build/_output/bin/kn-event-<OS>-<ARCH> \
    ~/.config/kn/plugins/kn-event
```

Check if plugin is loaded

```sh
$ kn -h
```

Run it

```sh
$ kn event -h
```

## Building

If you'd like to build the plugin yourself, you will need to have
[Golang](https://golang.org/) installed on your machine. Check the `go.mod` file
for current minimum version.

To build the plugin, just invoke the following script:

```sh
$ ./mage
```

Project will check, test and build binaries for various platforms. You'll get an
executable plugin binary named `kn-event-<OS>-<ARCH>` in `build/_output/bin`
dir.

You could list all available build targets with:

```sh
$ ./mage -l
```

### Updating dependencies

To update dependencies, please utilize the standard `hack/update-deps.sh`
script. It's also needed to run this script if you are doing any changes to Go
libraries. Read more about it at [knative/hack](https://github.com/knative/hack)
repository.
