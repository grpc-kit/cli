#!/usr/bin/env bash

set -euo pipefail

readonly PROTOC_VERSION="21.12"
readonly PROTOC_GEN_GO_VERSION="v1.36.11"
readonly PROTOC_GEN_GO_GRPC_VERSION="v1.6.2"
readonly GRPC_GATEWAY_VERSION="v2.30.0"

verify_sha256() {
  local checksum=$1 archive=$2
  if command -v sha256sum >/dev/null 2>&1; then
    printf '%s  %s\n' "${checksum}" "${archive}" | sha256sum -c -
    return
  fi
  if command -v shasum >/dev/null 2>&1; then
    printf '%s  %s\n' "${checksum}" "${archive}" | shasum -a 256 --check
    return
  fi
  echo "sha256sum or shasum is required to verify protoc" >&2
  return 1
}

install_protoc() {
  local host_os host_arch target_os target_arch archive checksum download_dir go_path
  host_os=$(go env GOHOSTOS)
  host_arch=$(go env GOHOSTARCH)
  target_os=${host_os}
  target_arch=${host_arch}

  case "${host_os}" in
    linux) ;;
    darwin) target_os="osx" ;;
    *) echo "unsupported protoc operating system: ${host_os}" >&2; return 1 ;;
  esac
  case "${host_arch}" in
    amd64) target_arch="x86_64" ;;
    arm64) target_arch="aarch_64" ;;
    *) echo "unsupported protoc architecture: ${host_arch}" >&2; return 1 ;;
  esac

  case "${target_os}/${target_arch}" in
    linux/x86_64) checksum="3a4c1e5f2516c639d3079b1586e703fc7bcfa2136d58bda24d1d54f949c315e8" ;;
    linux/aarch_64) checksum="2dd17f75d66a682640b136e31848da9fb2eefe68d55303baf8b32617374f6711" ;;
    osx/x86_64) checksum="9448ff40278504a7ae5139bb70c962acc78c32d8fc54b4890a55c14c68b9d10a" ;;
    osx/aarch_64) checksum="96839af0caed64352442fc8236f4bdf7c1cd6efcfaa98fa5db37307a73fc7c70" ;;
  esac

  archive="protoc-${PROTOC_VERSION}-${target_os}-${target_arch}.zip"
  download_dir=$(mktemp -d "${TMPDIR:-/tmp}/grpc-kit-protoc.XXXXXX")
  trap 'rm -rf "${download_dir}"' RETURN
  curl --fail --location --silent --show-error --retry 3 \
    "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/${archive}" \
    --output "${download_dir}/${archive}"
  verify_sha256 "${checksum}" "${download_dir}/${archive}"
  unzip -q "${download_dir}/${archive}" -d "${download_dir}/protoc"

  go_path=$(go env GOPATH)
  mkdir -p "${go_path}/bin" "${go_path}/include"
  install -m 0755 "${download_dir}/protoc/bin/protoc" "${go_path}/bin/protoc"
  rm -rf "${go_path}/include/google"
  cp -R "${download_dir}/protoc/include/google" "${go_path}/include/"
  rm -rf "${download_dir}"
  trap - RETURN
}

install_protoc_gen_go() {
  go install "google.golang.org/protobuf/cmd/protoc-gen-go@${PROTOC_GEN_GO_VERSION}"
}

install_protoc_gen_go_grpc() {
  go install "google.golang.org/grpc/cmd/protoc-gen-go-grpc@${PROTOC_GEN_GO_GRPC_VERSION}"
}

install_protoc_gen_grpc_gateway() {
  go install "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@${GRPC_GATEWAY_VERSION}"
}

install_protoc_gen_openapiv2() {
  go install "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@${GRPC_GATEWAY_VERSION}"
}

case "${1:-}" in
  protoc) install_protoc ;;
  protoc-gen-go) install_protoc_gen_go ;;
  protoc-gen-go-grpc) install_protoc_gen_go_grpc ;;
  protoc-gen-grpc-gateway) install_protoc_gen_grpc_gateway ;;
  protoc-gen-openapiv2) install_protoc_gen_openapiv2 ;;
  *)
    echo "usage: $0 {protoc|protoc-gen-go|protoc-gen-go-grpc|protoc-gen-grpc-gateway|protoc-gen-openapiv2}" >&2
    exit 2
    ;;
esac
