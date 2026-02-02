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

set -Eeuo pipefail

# Coordinates in GitHub.
ORG_NAME="${ORG_NAME:-knative-extensions}"

source "$(go run knative.dev/hack/cmd/script release.sh)"

PLUGIN="kn-event"

function build_release() {
  export GO111MODULE=on
  export CGO_ENABLED=0

  # Source build flags
  source "$(dirname "$0")/build-flags.sh"
  local ld_flags
  ld_flags="$(build_flags)"

  # Build the sender image with ko
  echo "🐳 Building kn-event-sender image"
  ko build --platform=linux/amd64,linux/arm64,linux/s390x,linux/ppc64le \
    --base-import-paths --tags="${TAG}" \
    ./cmd/kn-event-sender

  # Build CLI binaries for all platforms
  echo "🚧 🐧 Building for Linux (amd64)"
  GOOS=linux GOARCH=amd64 go build -ldflags "${ld_flags}" -o ./${PLUGIN}-linux-amd64 ./cmd/kn-event/...
  echo "🚧 💪 Building for Linux (arm64)"
  GOOS=linux GOARCH=arm64 go build -ldflags "${ld_flags}" -o ./${PLUGIN}-linux-arm64 ./cmd/kn-event/...
  echo "🚧 🐧 Building for Linux (ppc64le)"
  GOOS=linux GOARCH=ppc64le go build -ldflags "${ld_flags}" -o ./${PLUGIN}-linux-ppc64le ./cmd/kn-event/...
  echo "🚧 🐧 Building for Linux (s390x)"
  GOOS=linux GOARCH=s390x go build -ldflags "${ld_flags}" -o ./${PLUGIN}-linux-s390x ./cmd/kn-event/...
  echo "🚧 🍏 Building for macOS (amd64)"
  GOOS=darwin GOARCH=amd64 go build -ldflags "${ld_flags}" -o ./${PLUGIN}-darwin-amd64 ./cmd/kn-event/...
  echo "🚧 🍎 Building for macOS (arm64)"
  GOOS=darwin GOARCH=arm64 go build -ldflags "${ld_flags}" -o ./${PLUGIN}-darwin-arm64 ./cmd/kn-event/...
  echo "🚧 🎠 Building for Windows (amd64)"
  GOOS=windows GOARCH=amd64 go build -ldflags "${ld_flags}" -o ./${PLUGIN}-windows-amd64.exe ./cmd/kn-event/...

  ARTIFACTS_TO_PUBLISH="${PLUGIN}-linux-amd64 ${PLUGIN}-linux-arm64 ${PLUGIN}-linux-ppc64le ${PLUGIN}-linux-s390x ${PLUGIN}-darwin-amd64 ${PLUGIN}-darwin-arm64 ${PLUGIN}-windows-amd64.exe"

  sha256sum ${ARTIFACTS_TO_PUBLISH} > checksums.txt
  ARTIFACTS_TO_PUBLISH="${ARTIFACTS_TO_PUBLISH} checksums.txt"

  echo "🧮 Checksums:"
  cat checksums.txt
}

main "$@"
