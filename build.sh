#!/bin/sh
curl -O https://github.com/AlexKnauth/silksong-autosplit-wasm/releases/latest/download/silksong_autosplit_wasm_stable.wasm
curl -O https://raw.githubusercontent.com/LiveSplit/LiveSplit.AutoSplitters/master/LiveSplit.AutoSplitters.xml
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -H=windowsgui" -o sssplitmaker.exe github.com/CuteReimu/sssplitmaker
