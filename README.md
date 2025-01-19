# Elysium-project


sudo docker run --env-file local.env -d --name postgres -p 5432:5432 postgres:latest 


curl -X POST http://localhost:8080/peer -H "Content-Type: application/json" -d '{"public_key": "samplePublicKey", "OS_Arch": "x86_64-unknown-linux-musl"}'