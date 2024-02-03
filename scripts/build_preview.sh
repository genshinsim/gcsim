#!/bin/bash

rm -rf ./backend/cmd/preview/dist
cd ui
yarn
yarn build:embed
mv ./packages/embed/dist ../backend/cmd/preview/dist