#!/usr/bin/env bash

# Copyright 2021 The Knative Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o pipefail

# =================================================
# CUSTOMIZE ME:

# Name of the plugin
PLUGIN="kn-event"

# Directories containing go code which needs to be formatted
SOURCE_DIRS="cmd pkg internal"

# =================================================

# Store for later
if [ -z "$1" ]; then
    ARGS=("")
else
    ARGS=("$@")
fi

set -eu

# Run build
run() {
  # Switch on modules unconditionally
  export GO111MODULE=on

  # Jump into project directory
  pushd "$(basedir)" >/dev/null 2>&1

  # Print help if requested
  if $(has_flag --help -h); then
    display_help
    exit 0
  fi

  # Fast mode: Only compile
  if $(has_flag --fast -f); then
    go_build
    exit 0
  fi

  # Run only tests
  if $(has_flag --test -t); then
    go_test
    exit 0
  fi

  # Cross compile only
  if $(has_flag --all -x); then
    cross_build || { echo "Cross platform build failed" && exit 1; }
    exit 0
  fi

  # Default flow
  go_build
  go_test

  echo "────────────────────────────────────────────"
  ./$PLUGIN version
}

go_build() {
  echo "🚧 Compile"
  source "$(basedir)/hack/build-flags.sh"
  go build -ldflags "$(build_flags)" -o $PLUGIN ./cmd/kn-event/...
}

go_test() {
  local test_output
  test_output=$(mktemp /tmp/${PLUGIN}-output.XXXXXX)

  local red=""
  local reset=""
  # Use color only when a terminal is set
  if [ -t 1 ]; then
    red="[31m"
    reset="[39m"
  fi

  echo "🧪 Test"
  set +e
  go test -v ./pkg/... ./internal/... >$test_output 2>&1
  local err=$?
  if [ $err -ne 0 ]; then
    echo "🔥 ${red}Failure${reset}"
    cat $test_output | sed -e "s/^.*\(FAIL.*\)$/$red\1$reset/"
    rm $test_output
    exit $err
  fi
  rm $test_output
}

cross_build() {
  source "$(basedir)/hack/build-flags.sh"
  local ld_flags
  ld_flags="$(build_flags)"
  local failed=0

  echo "🔨 Cross-compile"

  export CGO_ENABLED=0

  echo "   🐧 ${PLUGIN}-linux-amd64"
  GOOS=linux GOARCH=amd64 go build -ldflags "${ld_flags}" -o ./${PLUGIN}-linux-amd64 ./cmd/kn-event/... || failed=1
  echo "   🐧 ${PLUGIN}-linux-arm64"
  GOOS=linux GOARCH=arm64 go build -ldflags "${ld_flags}" -o ./${PLUGIN}-linux-arm64 ./cmd/kn-event/... || failed=1
  echo "   🐧 ${PLUGIN}-linux-ppc64le"
  GOOS=linux GOARCH=ppc64le go build -ldflags "${ld_flags}" -o ./${PLUGIN}-linux-ppc64le ./cmd/kn-event/... || failed=1
  echo "   🐧 ${PLUGIN}-linux-s390x"
  GOOS=linux GOARCH=s390x go build -ldflags "${ld_flags}" -o ./${PLUGIN}-linux-s390x ./cmd/kn-event/... || failed=1
  echo "   🍏 ${PLUGIN}-darwin-amd64"
  GOOS=darwin GOARCH=amd64 go build -ldflags "${ld_flags}" -o ./${PLUGIN}-darwin-amd64 ./cmd/kn-event/... || failed=1
  echo "   🍎 ${PLUGIN}-darwin-arm64"
  GOOS=darwin GOARCH=arm64 go build -ldflags "${ld_flags}" -o ./${PLUGIN}-darwin-arm64 ./cmd/kn-event/... || failed=1
  echo "   🎠 ${PLUGIN}-windows-amd64.exe"
  GOOS=windows GOARCH=amd64 go build -ldflags "${ld_flags}" -o ./${PLUGIN}-windows-amd64.exe ./cmd/kn-event/... || failed=1

  return ${failed}
}

# Dir where this script is located
basedir() {
    # Default is current directory
    local script=${BASH_SOURCE[0]}

    # Resolve symbolic links
    if [ -L $script ]; then
        if readlink -f $script >/dev/null 2>&1; then
            script=$(readlink -f $script)
        elif readlink $script >/dev/null 2>&1; then
            script=$(readlink $script)
        elif realpath $script >/dev/null 2>&1; then
            script=$(realpath $script)
        else
            echo "ERROR: Cannot resolve symbolic link $script"
            exit 1
        fi
    fi

    local dir
    dir=$(dirname "$script")
    local full_dir
    full_dir=$(cd "${dir}/.." && pwd)
    echo "${full_dir}"
}

# Checks if a flag is present in the arguments.
has_flag() {
    filters="$@"
    for var in "${ARGS[@]}"; do
        for filter in $filters; do
          if [ "$var" = "$filter" ]; then
              echo 'true'
              return
          fi
        done
    done
    echo 'false'
}

# Display a help message.
display_help() {
    cat <<EOT
Build script for Kn plugin ${PLUGIN}

Usage: $(basename $BASH_SOURCE) [... options ...]

with the following options:

-f  --fast                    Only compile (without testing)
-t  --test                    Run tests only
-x  --all                     Build cross platform binaries
-h  --help                    Display this help message
    --verbose                 More output
    --debug                   Debug information for this script (set -x)

Examples:

* Compile and test: ..................... build.sh
* Compile only: ......................... build.sh --fast
* Run only tests: ....................... build.sh --test
* Build cross platform binaries: ........ build.sh --all
EOT
}

if $(has_flag --debug); then
    export PS4='+(${BASH_SOURCE[0]}:${LINENO}): ${FUNCNAME[0]:+${FUNCNAME[0]}(): }'
    set -x
fi

run "$@"
