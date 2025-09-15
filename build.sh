#!/bin/sh
curl -O https://github.com/AlexKnauth/silksong-autosplit-wasm/releases/download/0.1.3/silksong_autosplit_wasm_stable.wasm
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -H=windowsgui" -o sssplitmaker.exe github.com/CuteReimu/sssplitmaker
