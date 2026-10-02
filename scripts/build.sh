#!/bin/bash

set -eu

release_tag="$(git tag --sort=-version:refname | head -n 1)"

tags=()
if [ -n "${GCSIM_SHARE_KEY:-}" ]; then
	go run ./cmd/sharekeygen
	tags=(-tags sharekey)
fi

for name in gcsim server; do
	for os in darwin linux windows; do
		for arch in amd64 arm64; do
			out="${name}_${os}_${arch}"
			[ "${os}" != windows ] || out="${out}.exe"

			echo "building ${out}"
			CGO_ENABLED=0 GOOS="${os}" GOARCH="${arch}" go build \
				-trimpath \
				${tags[@]+"${tags[@]}"} \
				-ldflags "-X 'main.version=${release_tag}'" \
				-o "${out}" "./cmd/${name}"
		done
	done
done
