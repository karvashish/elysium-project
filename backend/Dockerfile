FROM golang:1.20-alpine AS build

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/. .
RUN go build -o main .

FROM alpine:latest

WORKDIR /root/

# temp
RUN apk add --no-cache wireguard-tools 

COPY --from=build /app/main .
COPY backend/migrations ./migrations
COPY .env /root/.env

EXPOSE 8080
CMD ["./main", "-env", "/root/.env"]
