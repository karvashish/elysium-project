#!/bin/sh

mkdir -p /certificates/root /certificates/intermediate

openssl genrsa -out /certificates/root/root.key 2048
openssl req -x509 -new -nodes -key /certificates/root/root.key -sha256 -days 1024 -out /certificates/root/root.crt -subj "/C=US/ST=State/L=City/O=Company/OU=Org/CN=RootCA"

openssl genrsa -out /certificates/intermediate/intermediate.key 2048
openssl req -new -key /certificates/intermediate/intermediate.key -out /certificates/intermediate/intermediate.csr -subj "/C=US/ST=State/L=City/O=Company/OU=Org/CN=IntermediateCA"

openssl x509 -req -in /certificates/intermediate/intermediate.csr -CA /certificates/root/root.crt -CAkey /certificates/root/root.key -CAcreateserial -out /certificates/intermediate/intermediate.crt -days 500 -sha256
