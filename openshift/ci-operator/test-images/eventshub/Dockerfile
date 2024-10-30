ARG GO_BUILDER=registry.ci.openshift.org/openshift/release:rhel-8-release-golang-1.22-openshift-4.17
ARG GO_RUNTIME=registry.access.redhat.com/ubi8/ubi-minimal

FROM $GO_BUILDER as builder

COPY . /code
WORKDIR /code

# Reset the goflags to avoid the -mod=vendor flag
ENV GOFLAGS=''

RUN go build \
    -o ./build/eventshub \
    knative.dev/reconciler-test/cmd/eventshub

FROM $GO_RUNTIME
COPY --from=builder /code/build/eventshub /usr/bin/eventshub
USER 65532
ENTRYPOINT ["/usr/bin/eventshub"]
