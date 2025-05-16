#!/bin/bash
set -e

readonly service="$1"
readonly output_dir="$2"
readonly package="$3"

oapi-codegen -generate types -o "$output_dir/openapi_types.gen.go" -package "$package" "api/openapi/dist/openapi.yaml"
oapi-codegen -generate strict-server,gin-server -o "$output_dir/openapi_api.gen.go" -package "$package" "api/openapi/dist/openapi.yaml"
oapi-codegen -generate types -o "src/go/pkg/api/genapi/client/openapi_types.gen.go" -package "$service" "api/openapi/dist/openapi.yaml"
oapi-codegen -generate client -o "src/go/pkg/api/genapi/client/openapi_client_gen.go" -package "$service" "api/openapi/dist/openapi.yaml"