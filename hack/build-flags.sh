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

function build_flags() {
  local version="${TAG:-}"
  # Use vYYYYMMDD-local-<hash> for the version string, if not passed.
  if [[ -z "${version}" ]]; then
    # Get the commit, excluding any tags but keeping the "dirty" flag
    local commit
    commit="$(git describe --always --dirty --match '^$')"
    [[ -n "${commit}" ]] || { echo "error getting the current commit" && exit 1; }
    version="v$(date +%Y%m%d)-local-${commit}"
  fi

  local image_basename="${IMAGE_BASENAME:-${KO_DOCKER_REPO:-}}"

  local pkg="knative.dev/kn-plugin-event/pkg/metadata"
  echo "-X '${pkg}.Version=${version}' -X '${pkg}.ImageBasename=${image_basename}'"
}
