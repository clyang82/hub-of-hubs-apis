FROM quay.io/bitnami/golang:1.17 AS builder
WORKDIR /go/src/github.com/clyang82/hub-of-hubs-apis
COPY . .

RUN make build --warn-undefined-variables

FROM registry.access.redhat.com/ubi8/ubi-minimal:latest
ENV USER_UID=10001

COPY --from=builder /go/src/github.com/clyang82/hub-of-hubs-apis/hub-of-hubs-apis /
RUN microdnf update && microdnf clean all

USER ${USER_UID}

