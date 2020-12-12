FROM golang:1.14.12-alpine3.11 as base
RUN apk add --update --no-cache make bash curl git build-base
ARG EXTRACT_PATH="/tmp/extract"
ENV PATH=/go/bin:$PATH
RUN mkdir -p ${EXTRACT_PATH}
ARG GO111MODULE=on
# Install bash-static
ARG BASH_STATIC_VERSION="5.0"
RUN curl -sL# https://github.com/robxu9/bash-static/releases/download/${BASH_STATIC_VERSION}/bash-linux -o ${EXTRACT_PATH}/bash-static \
  && chmod +x ${EXTRACT_PATH}/bash-static \
  && mv ${EXTRACT_PATH}/bash-static /usr/local/bash


# Build the manager binary
FROM base as builder
# Use nonroot user for building container & running CI (some tests assume running as non-root)
ENV USER dev
ENV HOME /home/dev
RUN adduser dev --disabled-password --home $HOME
RUN mkdir /workspace && chown dev:dev /workspace
USER dev
ENV GOPATH=/go
WORKDIR /workspace
# Copy the Go Modules manifests
COPY --chown=dev:dev go.mod go.mod
COPY --chown=dev:dev go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download
# Copy the go source
COPY --chown=dev:dev main.go main.go
COPY --chown=dev:dev pkg/ pkg/
COPY --chown=dev:dev cmd/ cmd/
# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o build/app main.go


# Make CI
FROM builder as ci
COPY Makefile Makefile
ARG GO111MODULE=on
# Fail if any files are not gofmt'd
RUN fmtFiles=$(go fmt ./...); test -z $fmtFiles || (echo "**** Need to run go fmt on: $fmtFiles ****" && exit 1)
# Install deps
RUN go install
RUN mcq setup
RUN mcq ci


# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot as final
WORKDIR /
COPY --chown=nonroot:nonroot --from=builder /workspace/build/app .
COPY --chown=nonroot:nonroot --from=builder /workspace/build/app .
COPY --chown=nonroot:nonroot --from=base /usr/local/bash .
USER nonroot:nonroot
ENTRYPOINT ["/app"]
