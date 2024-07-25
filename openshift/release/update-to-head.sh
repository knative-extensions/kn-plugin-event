#!/usr/bin/env bash

set -Eeuo pipefail

project_dir="$(realpath "$(dirname "${BASH_SOURCE[0]:-$0}")/../..")"

cd "$project_dir/openshift"

exec go run \
  github.com/cardil/deviate/cmd/deviate \
  sync --config "${project_dir}/.deviate.yaml"
