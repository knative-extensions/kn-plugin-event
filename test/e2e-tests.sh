#!/usr/bin/env bash

set -Eeo pipefail

# shellcheck disable=SC1090
source "$(go run knative.dev/hack/cmd/script e2e-tests.sh)"

function start_latest_knative_serving() {
  local KNATIVE_NET_KOURIER_RELEASE
  KNATIVE_NET_KOURIER_RELEASE="$(get_latest_knative_yaml_source "net-kourier" "kourier")"
  start_knative_serving "${KNATIVE_SERVING_RELEASE_CRDS}" \
    "${KNATIVE_SERVING_RELEASE_CORE}" \
    "${KNATIVE_NET_KOURIER_RELEASE}"

  kubectl patch configmap/config-network \
    --namespace knative-serving \
    --type merge \
    --patch '{"data":{"ingress.class":"kourier.ingress.networking.knative.dev"}}'
}

function knative_setup() {
  start_latest_knative_serving
  start_latest_knative_eventing
}

initialize "$@"

set -Eeuo pipefail

./hack/publish.sh
./hack/build.sh --fast

export KN_PLUGIN_EVENT_EXECUTABLE="${PWD}/kn-event"

go_test_e2e -timeout 10m ./test/... || fail_test 'kn-event e2e tests'

success
