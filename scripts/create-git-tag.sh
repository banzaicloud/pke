#!/bin/bash -euf

OWNER="$1"
REPO="$2"
REF="$3"
TAG="$4"

API_URL="https://api.github.com/repos/${OWNER}/${REPO}"

SHA=$(curl -sS -X GET -H "Authorization: token ${GITHUB_TOKEN}" "${API_URL}/git/${REF}" | jq -r '.object.sha')

curl \
    -sS \
    -X POST \
    -H "Authorization: token ${GITHUB_TOKEN}" \
    -H "Content-Type: application/json" \
    -d '{"sha": "'"${SHA}"'", "ref": "refs/tags/'"${TAG}"'"}' \
    "${API_URL}/git/refs"
