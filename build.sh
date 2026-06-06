#!/bin/sh
set -e
curl -L -O https://github.com/AlexKnauth/silksong-autosplit-wasm/releases/latest/download/silksong_autosplit_wasm_stable.wasm
curl -O https://raw.githubusercontent.com/LiveSplit/LiveSplit.AutoSplitters/master/LiveSplit.AutoSplitters.xml
wails build
wails build -platform=windows/amd64 -webview2 embed
