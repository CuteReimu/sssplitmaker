package main

import (
	_ "embed"
	"github.com/lxn/walk"
	"os"
	"path/filepath"
)

//go:embed silksong_autosplit_wasm_stable.wasm
var wasmFile []byte

func saveWasmFile() {
	dlg := new(walk.FileDialog)
	dlg.Title = "保存wasm文件"
	dlg.Filter = "wasm文件（*.wasm）|*.wasm"
	dlg.FilePath = "silksong_autosplit_wasm_stable.wasm"
	if ok, err := dlg.ShowSave(mainWindow); err != nil {
		walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
	} else if ok {
		file := dlg.FilePath
		if filepath.Ext(file) != ".wasm" {
			file += ".wasm"
		}
		err = os.WriteFile(file, wasmFile, 0644)
		if err != nil {
			walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
			return
		}
	}
}
