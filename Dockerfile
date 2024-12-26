# syntax=docker/dockerfile:1

ARG GO_VERSION=1.23.4
ARG ALPINE_VERSION=3.21

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

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/fetch_task ./cmd/fetch_task

FROM alpine:${ALPINE_VERSION} AS final

RUN --mount=type=cache,target=/var/cache/apk \
    apk --update add \
        ca-certificates \
        tzdata \
        && \
        update-ca-certificates

ADD --chmod=111 'https://github.com/apple/pkl/releases/download/0.27.1/pkl-alpine-linux-amd64' /bin/pkl

ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    tikhon

RUN mkdir pkl_cache && chown -R tikhon:tikhon /pkl_cache && chmod 755 /pkl_cache

ARG SOURCE_COMMIT
RUN echo $SOURCE_COMMIT > release.txt && chown tikhon release.txt && chmod 755 release.txt

USER tikhon

COPY --from=build /bin/server /bin/
COPY --from=build /bin/fetch_task /bin/
COPY . .


EXPOSE 9990
