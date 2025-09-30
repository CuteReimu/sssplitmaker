package main

import (
	"github.com/CuteReimu/sssplitmaker/splitmaker"
	"os"
	"os/exec"
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
var commentTextLabel *walk.TextLabel
var startTriggerComboBox *walk.ComboBox
var splitmakerComboBox *walk.ComboBox
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
					TextLabel{TextAlignment: AlignHNearVCenter, Text: "或者把文件拖拽进来"},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					TextLabel{
						TextAlignment: AlignHFarVCenter,
						Text:          "一键填入预设分割点(请配合+和-使用)",
					},
					ComboBox{
						AssignTo: &splitmakerComboBox,
						Model:    splitmaker.GetAllFiles(),
						Value:    "",
						Enabled:  false,
						OnCurrentIndexChanged: func() {
							loadLayoutFileFromSplitmaker(splitmakerComboBox.Text())
						},
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					TextLabel{
						TextAlignment: AlignHFarVCenter,
						Text:          "Auto Splitter Version: 0.1.13",
					},
					PushButton{
						Text:      "获取wasm文件",
						Alignment: AlignHFarVCenter,
						OnClicked: saveWasmFile,
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					TextLabel{
						TextAlignment: AlignHFarVCenter,
						Text:          "由于翻译水平有限，翻译可能不准确，请检查翻译对照列表是否正确后使用。",
					},
					PushButton{
						Text:      "打开翻译对照列表",
						Alignment: AlignHFarVCenter,
						OnClicked: func() {
							const url = "https://cutereimu.cn/daily/silksong/splits-translate.html"
							if err := exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start(); err != nil {
								walk.MsgBox(mainWindow, "错误", "打开浏览器失败，请手动访问："+url, walk.MsgBoxIconError)
							}
						},
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
					TextLabel{
						AssignTo:      &commentTextLabel,
						TextAlignment: AlignHNearVCenter,
						Text:          "想要修改左边一列，请直接在LiveSplit中使用Edit Splits进行修改。",
						Visible:       false,
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
					PushButton{
						MaxSize:   Size{Width: 100},
						Alignment: AlignHFarVCenter,
						Text:      "帮助",
						OnClicked: func() {
							const url = "https://cutereimu.cn/daily/silksong/sssplitmaker-faq.html"
							if err := exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start(); err != nil {
								walk.MsgBox(mainWindow, "错误", "打开浏览器失败，请手动访问："+url, walk.MsgBoxIconError)
							}
						},
					},
					PushButton{Text: "清空", OnClicked: func() {
						resetLines(1)
						fileLayoutData = nil
						fileWasmSettings = nil
						commentTextLabel.SetVisible(false)
						splitmakerComboBox.SetEnabled(false)
						saveButton.SetEnabled(false)
						err := saveButton.SetText("请先打开lsl文件")
						if err != nil {
							walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
						}
					}},
					PushButton{AssignTo: &saveButton, Text: "请先打开lsl文件", OnClicked: onSaveLayoutFile, Enabled: false},
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
