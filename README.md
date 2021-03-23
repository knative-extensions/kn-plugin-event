# kn-plugin-event

`kn-plugin-event` is a plugin of Knative Client, for managing cloud events from
command line.

## Description

With this plugin, you can build and send the cloud events to publicly available
addresses (URLs) or to Addressable (Kubernetes service, Knative service, broker,
channels, etc).

## Build and Install

You must
[set up your development environment](https://github.com/knative/client/blob/master/docs/DEVELOPMENT.md#prerequisites)
before you build.

**Building:**

Once you've set up your development environment, let's build the plugin.

```sh
$ go build -o kn-event ./cmd/kn-event/main.go
```

You'll get an executable plugin binary namely `kn-event` in your current dir.
You're ready to use `kn-event` as a stand alone binary, check the available
commands `./kn-event -h`.

**Installing:**

If you'd like to use the plugin with `kn` CLI, install the plugin by simply
copying the executable file under `kn` plugins directory as:

```sh
mkdir -p ~/.config/kn/plugins
cp kn-event ~/.config/kn/plugins
```

Check if plugin is loaded

```sh
kn -h
```

Run it

```sh
kn event -h
```

## Examples

**Send an event:**

Send an event to a Knative service `event-display` in namespace `default`:

```sh
$ kn event send \
    --type org.example.kn.ping \
    --id $(uuidgen) \
    --field event.type=test \
    --field event.data=ping \
    --to Service:v1:event-display \
    --namespace default \
```
