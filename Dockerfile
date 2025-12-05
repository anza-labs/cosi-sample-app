# Easy crosscomple toolkit
FROM --platform=$BUILDPLATFORM tonistiigi/xx:1.9.0 AS xx

# Build the cosi-sample-app binary
FROM --platform=$BUILDPLATFORM docker.io/library/golang:1.25 AS builder
ARG TARGETOS
ARG TARGETARCH
ARG TARGETPLATFORM
COPY --from=xx / /

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN xx-go mod download

# Copy the go source
COPY main.go main.go
COPY pkg/ pkg/

# Build
ENV CGO_ENABLED=0
RUN xx-go build -trimpath -a -o cosi-sample-app main.go && \
    xx-verify cosi-sample-app

# Use distroless as minimal base image to package the cosi-sample-app binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/cosi-sample-app .
USER 65532:65532

ENTRYPOINT ["/cosi-sample-app"]
