#!/bin/bash

rm -rf ./backend/cmd/preview/dist
cd ui
pnpm install
pnpm build:embed
mv ./packages/embed/dist ../backend/cmd/preview/dist