Invoke-WebRequest -Uri 'https://github.com/AlexKnauth/silksong-autosplit-wasm/releases/latest/download/silksong_autosplit_wasm_stable.wasm' -OutFile 'silksong_autosplit_wasm_stable.wasm' -v
Invoke-WebRequest -Uri 'https://raw.githubusercontent.com/LiveSplit/LiveSplit.AutoSplitters/master/LiveSplit.AutoSplitters.xml' -OutFile 'LiveSplit.AutoSplitters.xml' -v
wails build -webview2 embed
wails build -platform=darwin/arm64