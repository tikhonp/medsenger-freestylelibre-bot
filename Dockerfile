# syntax=docker/dockerfile:1

ARG GOVERSION=1.24.3
ARG PKL_VERSION=0.28.2

FROM golang:${GOVERSION} AS dev
RUN go install "github.com/air-verse/air@latest"
RUN go install "github.com/apple/pkl-go/cmd/pkl-gen-go@latest"
RUN go install "github.com/a-h/templ/cmd/templ@latest"
ARG PKL_VERSION
RUN curl -L -o /usr/bin/pkl "https://github.com/apple/pkl/releases/download/${PKL_VERSION}/pkl-linux-$(uname -m)" && chmod +x /usr/bin/pkl
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download && go mod verify
CMD ["air", "-c", ".air.toml"]


FROM golang:${GOVERSION} AS prod
ARG TARGETARCH
ARG PKL_VERSION
ADD --chmod=111 "https://github.com/apple/pkl/releases/download/${PKL_VERSION}/pkl-linux-${TARGETARCH}" /bin/pkl
WORKDIR /src
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    GOARCH=$TARGETARCH go build -o /bin/server ./cmd/server
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    GOARCH=$TARGETARCH go build -o /bin/fetch_task ./cmd/fetch_task
COPY . .
