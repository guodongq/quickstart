#!/bin/bash
set -e

#readonly service="$1"
#readonly output_dir="$2"
#readonly package="$3"

readonly output_dir="$1"
readonly package="$2"

#oapi-codegen -generate="server,strict-server,spec" -o "$output_dir/$service/openapi_api.gen.go" -package "$service" "api/openapi/dist/$service.openapi.yaml"
#oapi-codegen -generate="types" -o "$output_dir/$service/openapi_types.gen.go" -package "$service" "api/openapi/dist/$service.openapi.yaml"
#oapi-codegen -generate="client" -o "$output_dir/$service/openapi_client_gen.go" -package "$service" "api/openapi/dist/$service.openapi.yaml"

# std-http
oapi-codegen -generate "gin-server,strict-server, spec" -o "$output_dir/openapi_server.gen.go" -package "$package" "api/openapi/dist/openapi.yaml"
oapi-codegen -generate "types" -o "$output_dir/openapi_types.gen.go" -package "$package" "api/openapi/dist/openapi.yaml"
oapi-codegen -generate "client" -o "$output_dir/openapi_client_gen.go" -package "$package" "api/openapi/dist/openapi.yaml"