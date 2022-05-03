#!/bin/bash

set -e

protoc_version="3.17.3"
protoc_gen_go_version="v1.26.0"
protoc_gen_go_grpc_version="1.1.0"

check_protoc_version() {
  if ! command -v protoc >/dev/null; then
    echo "command protoc could not be found"
    echo "please install protoc first. Ref: https://github.com/protocolbuffers/protobuf/releases"
    exit
  fi
  if [ ! -d "/usr/local/include/google/" ]; then
    echo "missing google protobuf"
    echo "please install it first. Ref: https://github.com/grpc-ecosystem/grpc-gateway/issues/422#issuecomment-409809309"
    exit
  fi
  version=$(protoc --version)
  if ! [[ ${version} == *"${protoc_version}"* ]]; then
    echo "invalid proto version ${version}"
    exit 1
  fi
}

check_protoc_gen_go_version() {
  version=$(protoc-gen-go --version)
  if ! [[ ${version} == *"${protoc_gen_go_version}"* ]]; then
    echo "invalid proto-gen-go version ${version}"
    exit 1
  fi
}

check_protoc_gen_go_grpc_version() {
  version=$(protoc-gen-go-grpc --version)
  if ! [[ ${version} == *"${protoc_gen_go_grpc_version}"* ]]; then
    echo "invalid proto-gen-go-grpc version ${version}"
    exit 1
  fi
}

check_protoc_gen_grpc_gateway() {
  if ! command -v protoc-gen-grpc-gateway &>/dev/null; then
    echo "command protoc-gen-grpc-gateway could not be found"
    echo "run 'go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.4.0' to install"
    exit 1
  fi
}

check_protoc_gen_openapiv2() {
  if ! command -v protoc-gen-openapiv2 &>/dev/null; then
    echo "command protoc-gen-openapiv2 could not be found"
    echo "run 'go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.4.0' to install"
    exit 1
  fi
}

generate_code() {
  # list all proto packages except \`google\` and \`protoc-gen-openapiv2\`
  for dir in $(find proto -mindepth 1 -maxdepth 1 -type d | egrep -v 'google' | egrep -v 'protoc-gen-openapiv2'); do
    echo $dir
    dir_name=$(basename ${dir})
    protoc -I=./proto \
      --go_opt=paths=source_relative \
      --go_out=./server/pb \
      --go-grpc_opt=paths=source_relative \
      --go-grpc_out=./server/pb \
      --grpc-gateway_opt=paths=source_relative \
      --grpc-gateway_out=./server/pb \
      --openapiv2_out ./server/openapi \
      proto/${dir_name}/*.proto
  done
}

cleanup() {
  rm -rf server/openapi/*
  rm -rf server/pb/*
}

check_protoc_version
check_protoc_gen_go_version
check_protoc_gen_go_grpc_version
check_protoc_gen_grpc_gateway
check_protoc_gen_openapiv2

cleanup
generate_code
