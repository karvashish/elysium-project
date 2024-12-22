#!/bin/bash

set -e

echo "Cleaning up the previous target directory..."
rm -rf ./target
mkdir -p ./target
echo "Target directory created."

echo "Copying backend Go modules and source files to target directory..."
mkdir -p ./target/backend
cp backend/go.mod backend/go.sum ./target/backend/
cp -r backend/ ./target/
cp -r backend/migrations ./target/backend/migrations
echo "Backend files copied."

(
  cd ./target/backend
  go build -o main .
) &

for i in {1..10}; do
  echo -ne "\rGo application building: /" && sleep 0.2
  echo -ne "\rGo application building: -" && sleep 0.2
  echo -ne "\rGo application building: \\" && sleep 0.2
  echo -ne "\rGo application building: |" && sleep 0.2
done

wait
echo -ne "\rGo application building: Done\n"


echo "Copying compiled binary and necessary files to the target directory..."
cp ./target/backend/main ./target/
cp -r backend/migrations ./target/migrations/
cp .env ./target/.env
cp -r client ./target/client/
echo "Files copied to the target directory."

echo "Cleaning up unnecessary backend directory in target..."
rm -rf ./target/backend
echo "Cleanup completed."
