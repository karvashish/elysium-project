FROM golang:1.23.2-alpine AS build

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/. .
RUN go build -o main .

FROM alpine:latest

WORKDIR /root/

RUN apk add --no-cache wireguard-tools cargo

COPY --from=build /app/main .
COPY backend/migrations ./migrations
COPY .env /root/.env
COPY client /root/client

EXPOSE 8080
CMD ["./main", "-env", "/root/.env"]
