FROM alpine:latest

RUN apk add --no-cache openssl

WORKDIR /certificates

COPY create_certs.sh /create_certs.sh
RUN chmod +x /create_certs.sh

CMD ["/create_certs.sh"]
