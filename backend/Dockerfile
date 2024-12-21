FROM golang:1.23.4-alpine AS build

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/. .
RUN go build -o main .

FROM alpine:latest

WORKDIR /root/

RUN apk add --no-cache \
    build-base \
    musl-dev \
    rustup \
    curl \
    binutils \
    gcc \
    g++ && \
    curl https://sh.rustup.rs -sSf | sh -s -- -y && \
    source $HOME/.cargo/env && \
    rustup install stable && \
    rustup default stable && \
    rustup target add x86_64-pc-windows-gnu x86_64-unknown-linux-musl aarch64-unknown-linux-musl && \
    rm -rf /var/cache/apk/* /tmp/* /var/tmp/*

ENV PATH="/root/.cargo/bin:${PATH}"

COPY --from=build /app/main .
COPY backend/migrations ./migrations
COPY .env /root/.env
COPY client /root/client

EXPOSE 8080
CMD ["./main", "-env", "/root/.env"]
