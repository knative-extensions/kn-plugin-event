FROM registry.ci.openshift.org/openshift/release:rhel-8-release-golang-1.21-openshift-4.16

# Install the oc and kubectl clients
RUN curl -LO https://mirror.openshift.com/pub/openshift-v4/x86_64/clients/ocp-dev-preview/candidate-4.16/openshift-client-linux-amd64-rhel8.tar.gz \
    && tar -xzf openshift-client-linux-amd64-rhel8.tar.gz \
    && chmod +x oc kubectl \
    && mv oc kubectl /usr/local/bin/

# Reset the goflags to avoid the -mod=vendor flag
ENV GOFLAGS=''
