#!/bin/sh

GOOS=js GOARCH=wasm go build -o ../server/dist/wasm.wasm