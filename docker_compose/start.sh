#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.


set -x

docker compose -f docker_compose/docker-compose.yml down --volumes || true
rm -f docker_compose/token.txt || true
docker compose -f docker_compose/docker-compose.yml up -d
until curl -s -X GET "http://localhost:4000/api/ready" | grep -q '"message":"ready"'; do
    echo "Waiting for the service to be ready..."
    sleep 5
done
echo "Service is ready!"
