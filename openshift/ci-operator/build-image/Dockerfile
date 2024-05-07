ARG OCP_VERSION=4.16
ARG GOLANG_VERSION=1.21
ARG RHEL_VERSION=8
FROM registry.ci.openshift.org/openshift/release:rhel-${RHEL_VERSION}-release-golang-${GOLANG_VERSION}-openshift-${OCP_VERSION}
ARG OCP_VERSION

# Install the oc and kubectl clients
RUN curl -LO https://mirror.openshift.com/pub/openshift-v4/x86_64/clients/ocp-dev-preview/candidate-${OCP_VERSION}/openshift-client-linux.tar.gz \
    && tar -xzf openshift-client-linux.tar.gz \
    && chmod +x oc kubectl \
    && mv oc kubectl /usr/local/bin/

# Reset the goflags to avoid the -mod=vendor flag
ENV GOFLAGS=''
