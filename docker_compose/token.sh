#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.

set -x

if [[ ! -f docker_compose/token.txt ]]; then
    expires=$(date -d "+1 day" -u +"%Y-%m-%dT%H:%M:%S.%3NZ")
    curl -s -X POST localhost:8080/api-keys -H 'Content-Type: application/json' -d "{\"name\":\"terraform_integration_test\",\"expiredAt\":\"${expires}\",\"allowedPath\":\".*\",\"type\":\"devlake\"}" | jq -r '.apiKey' > docker_compose/token.txt
fi

cat docker_compose/token.txt
