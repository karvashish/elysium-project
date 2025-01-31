FROM golang:1.23.5-alpine AS build

WORKDIR /app

COPY backend/go.mod backend/go.sum ./ 
RUN go mod download

COPY backend/. . 
RUN go build -o main .

FROM rust:1.84.1-slim AS runtime

WORKDIR /root/

RUN apt-get update && apt-get install -y --no-install-recommends musl-dev gcc g++ musl-tools build-essential gcc-x86-64-linux-gnu mingw-w64\
    && rustup target add x86_64-pc-windows-gnu x86_64-unknown-linux-musl aarch64-unknown-linux-musl \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ENV PATH="/root/.cargo/bin:${PATH}"

COPY --from=build /app/main . 
COPY backend/migrations ./migrations 
COPY .env /root/.env 
COPY client /root/client

EXPOSE 8080
CMD ["./main", "-env", "/root/.env"]
