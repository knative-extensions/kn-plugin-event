ARG OCP_VERSION=4.16
ARG GOLANG_VERSION=1.21
ARG RHEL_VERSION=8
FROM registry.ci.openshift.org/ocp/${OCP_VERSION}:tools as tools

FROM registry.ci.openshift.org/openshift/release:rhel-${RHEL_VERSION}-release-golang-${GOLANG_VERSION}-openshift-${OCP_VERSION}

COPY --from=tools /usr/bin/oc /usr/bin/
RUN ln -s /usr/bin/oc /usr/bin/kubectl

# Reset the goflags to avoid the -mod=vendor flag
ENV GOFLAGS=''
