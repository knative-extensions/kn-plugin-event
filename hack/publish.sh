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

# Publishes the kn-event-sender image to the container registry.
# Used by e2e tests and CI.

set -Eeuo pipefail

cd "$(dirname "$0")/.."

# Compute version tag (same logic as build-flags.sh)
version="${TAG:-}"
if [[ -z "${version}" ]]; then
  commit="$(git describe --always --dirty --match '^$')"
  version="v$(date +%Y%m%d)-local-${commit}"
fi

echo "🐳 Publishing kn-event-sender image with tag ${version}"
ko build --base-import-paths --tags="${version}" ./cmd/kn-event-sender
