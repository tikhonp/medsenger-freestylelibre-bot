# syntax=docker/dockerfile:1

ARG GO_VERSION=1.23.0
ARG ALPINE_VERSION=3.20.2

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build

WORKDIR /src

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

ARG TARGETARCH

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/server ./cmd/server


FROM alpine:${ALPINE_VERSION} AS final

RUN --mount=type=cache,target=/var/cache/apk \
    apk --update add \
        ca-certificates \
        tzdata \
        && \
        update-ca-certificates

ADD 'https://github.com/apple/pkl/releases/download/0.26.3/pkl-alpine-linux-amd64' /bin/pkl
RUN chmod +x /bin/pkl

ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    tikhon
USER tikhon

COPY --from=build /bin/server /bin/
COPY pkl pkl

EXPOSE 9990

ENTRYPOINT [ "/bin/server" ]
