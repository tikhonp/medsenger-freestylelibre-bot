# syntax=docker/dockerfile:1

ARG GOVERSION=1.24.5

FROM golang:${GOVERSION}-alpine AS dev
RUN go install "github.com/air-verse/air@latest" && go install "github.com/a-h/templ/cmd/templ@latest"
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download && go mod verify
CMD ["air", "-c", ".air.toml"]


FROM golang:${GOVERSION}-alpine AS build-prod
ARG TARGETARCH
WORKDIR /src
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/server ./cmd/server
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/fetch_task ./cmd/fetch_task

FROM alpine AS prod
WORKDIR /src
COPY --from=build-prod /bin/server /bin/fetch_task /bin/
COPY . .
