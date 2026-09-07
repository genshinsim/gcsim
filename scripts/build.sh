#!/bin/bash

set -eu

# shellcheck disable=SC2154
: "${GCSIM_SHARE_KEY}"

release_tag="$(git tag --sort=-version:refname | head -n 1)"

for name in gcsim server; do
	for os in darwin linux windows; do
		for arch in amd64 arm64; do
			out="${name}_${os}_${arch}"
			[ "${os}" != windows ] || out="${out}.exe"

			echo "building ${out}"
			CGO_ENABLED=0 GOOS="${os}" GOARCH="${arch}" go build \
				-trimpath \
				-ldflags "-X 'main.shareKey=${GCSIM_SHARE_KEY}' -X 'main.version=${release_tag}'" \
				-o "${out}" "./cmd/${name}"
		done
	done
done
