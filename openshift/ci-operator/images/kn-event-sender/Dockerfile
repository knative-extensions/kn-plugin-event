ARG GO_BUILDER=registry.ci.openshift.org/openshift/release:rhel-8-release-golang-1.22-openshift-4.17
ARG GO_RUNTIME=registry.access.redhat.com/ubi8/ubi-minimal

FROM $GO_BUILDER as builder

COPY . /code
WORKDIR /code

# Reset the goflags to avoid the -mod=vendor flag
ENV GOFLAGS=''
RUN go build \
    -o ./build/kn-event-sender \
    ./cmd/kn-event-sender

FROM $GO_RUNTIME
COPY --from=builder /code/build/kn-event-sender /usr/bin/kn-event-sender
USER 65532
ENTRYPOINT ["/usr/bin/kn-event-sender"]
