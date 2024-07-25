# kn-plugin-event

This repository holds OpenShift's fork of 
[`knative-sandbox/kn-plugin-event`](https://github.com/knative-sandbox/kn-plugin-event) 
with additions and fixes needed only for the OpenShift side of things.

[![Mage](https://github.com/cardil/kn-plugin-event-fork/actions/workflows/go.yml/badge.svg?branch=release-next)](https://github.com/cardil/kn-plugin-event-fork/actions/workflows/go.yml)
[![GolangCI-Lint](https://github.com/cardil/kn-plugin-event-fork/actions/workflows/golangci-lint.yaml/badge.svg?branch=release-next)](https://github.com/cardil/kn-plugin-event-fork/actions/workflows/golangci-lint.yaml)

## How this repository works ?

The `main` branch holds up-to-date specific [openshift files](./openshift)
that are necessary for CI setups and maintaining it. This includes:

- Scripts to create a new release branch from `upstream`
- CI setup files & tests scripts

Each release branch holds the upstream code for that release and our
OpenShift's specific files.

## CI Setup

For the CI setup, two repositories are of importance:

- This repository
- [openshift/release](https://github.com/openshift/release) which
  contains the configuration of CI jobs that are run on this
  repository
