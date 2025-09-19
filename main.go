package main

import (
	"os"
	"path/filepath"
	"syscall"

	"github.com/CuteReimu/sssplitmaker/translate"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

func getSystemMetrics(nIndex int) int {
	ret, _, _ := syscall.NewLazyDLL(`User32.dll`).NewProc(`GetSystemMetrics`).Call(uintptr(nIndex))
	return int(ret)
}

var mainWindow *walk.MainWindow
var splitLinesViewContainer *walk.Composite
var splitLinesView *walk.Composite
var startTriggerCheckBox *walk.CheckBox
var startTriggerComboBox *walk.ComboBox
var saveButton *walk.PushButton

func main() {
	screenX, screenY := getSystemMetrics(0), getSystemMetrics(1)
	width, height := 720, 960
	err := MainWindow{
		OnDropFiles: func(f []string) {
			if len(f) > 0 {
				file := f[0]
				if filepath.Ext(file) == ".lss" {
					buf, err := os.ReadFile(file)
					if err != nil {
						walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
						return
					}
					loadSplitFile(buf)
				} else if filepath.Ext(file) == ".lsl" {
					buf, err := os.ReadFile(file)
					if err != nil {
						walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
						return
					}
					loadLayoutFile(buf)
				}
			}
		},
		AssignTo: &mainWindow,
		Title:    "计时器生成器",
		Bounds:   Rectangle{X: (screenX - width) / 2, Y: (screenY - height) / 2, Width: width, Height: height},
		Layout:   VBox{},
		Children: []Widget{
			Composite{
				MaxSize: Size{Height: 20},
				Layout:  HBox{},
				Children: []Widget{
					TextLabel{TextAlignment: AlignHFarVCenter, Text: "你可以"},
					PushButton{Text: "打开lss文件", OnClicked: onClickLoadSplitFile},
					PushButton{Text: "打开lsl文件", OnClicked: onClickLoadLayoutFile},
					TextLabel{TextAlignment: AlignHFarVCenter, Text: "或者把文件拖拽进来"},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					TextLabel{
						TextAlignment: AlignHFarVCenter,
						Text:          "Auto Splitter Version: 0.1.7",
					},
					PushButton{
						Text:      "获取wasm文件",
						Alignment: AlignHFarVCenter,
						OnClicked: saveWasmFile,
					},
				},
			},
			ScrollView{
				HorizontalFixed: true,
				Layout:          VBox{},
				Children: []Widget{
					Composite{
						MaxSize: Size{Width: 0, Height: 25},
						Layout:  HBox{},
						Children: []Widget{
							CheckBox{
								AssignTo: &startTriggerCheckBox,
								Text:     "自动开始",
								Enabled:  false,
								Checked:  true,
							},
							ComboBox{
								AssignTo: &startTriggerComboBox,
								Model:    splitDescriptions,
								Value:    translate.GetSplitDescriptionByID("StartNewGame"),
							},
						},
					},
					Composite{
						AssignTo: &splitLinesViewContainer,
						Layout:   Flow{},
						Children: []Widget{
							Composite{
								AssignTo:  &splitLinesView,
								Alignment: AlignHCenterVNear,
								Layout:    VBox{},
							},
						},
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{Text: "清空", OnClicked: func() {
						resetLines(1)
						fileLayoutData = nil
						fileWasmSettings = nil
						saveButton.SetEnabled(false)
					}},
					PushButton{AssignTo: &saveButton, Text: "另存为", OnClicked: onSaveLayoutFile, Enabled: false},
				},
			},
		},
	}.Create()
	addLine()
	if err != nil {
		walk.MsgBox(nil, "错误", err.Error(), walk.MsgBoxIconError)
		return
	}
	hWnd := mainWindow.Handle()
	currStyle := win.GetWindowLong(hWnd, win.GWL_STYLE)
	win.SetWindowLong(hWnd, win.GWL_STYLE, currStyle & ^win.WS_SIZEBOX)
	mainWindow.Run()
}
