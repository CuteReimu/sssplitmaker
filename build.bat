@echo off
bitsadmin /transfer n https://github.com/AlexKnauth/silksong-autosplit-wasm/releases/download/0.1.7/silksong_autosplit_wasm_stable.wasm %~dp0\silksong_autosplit_wasm_stable.wasm
go build -ldflags "-s -w -H=windowsgui" -o sssplitmaker.exe github.com/CuteReimu/sssplitmaker
