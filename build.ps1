Invoke-WebRequest -Uri 'https://github.com/AlexKnauth/silksong-autosplit-wasm/releases/latest/download/silksong_autosplit_wasm_stable.wasm' -OutFile 'silksong_autosplit_wasm_stable.wasm' -v
Invoke-WebRequest -Uri 'https://raw.githubusercontent.com/LiveSplit/LiveSplit.AutoSplitters/master/LiveSplit.AutoSplitters.xml' -OutFile 'LiveSplit.AutoSplitters.xml' -v
go build -ldflags "-s -w -H=windowsgui" -o sssplitmaker.exe github.com/CuteReimu/sssplitmaker
go env -w GOOS=darwin
go env -w GOARCH=arm64
go build -ldflags "-s -w" -o sssplitmaker github.com/CuteReimu/sssplitmaker