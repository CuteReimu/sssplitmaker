package main

import (
	"context"
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
var startTriggerComboBox *walk.ComboBox
var splitmakerComboBox *walk.ComboBox

func main() {
	screenX, screenY := getSystemMetrics(0), getSystemMetrics(1)
	width, height := 720, 960
	err := MainWindow{
		OnDropFiles: func(f []string) {
			if len(f) > 0 {
				file := f[0]
				if filepath.Ext(file) != ".lss" {
					return
				}
				buf, err := os.ReadFile(file)
				if err != nil {
					walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
					return
				}
				loadSplitFile(buf)
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
					TextLabel{TextAlignment: AlignHNearVCenter, Text: "或者把文件拖拽进来，也可以使用现有模板"},
					ComboBox{
						AssignTo: &splitmakerComboBox,
						Model:    GetAllFiles(),
						Value:    "",
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
						Text:          "Auto Splitter Version: 1.13.0",
					},
					PushButton{
						Text:      "更新LiveSplit",
						Alignment: AlignHFarVCenter,
						OnClicked: fixLiveSplit,
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
							if err := exec.CommandContext(context.Background(), "rundll32", "url.dll,FileProtocolHandler", url).Start(); err != nil {
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
								Text:    "自动开始",
								Enabled: false,
								Checked: true,
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
					PushButton{
						MaxSize:   Size{Width: 100},
						Alignment: AlignHFarVCenter,
						Text:      "帮助",
						OnClicked: func() {
							const url = "https://cutereimu.cn/daily/silksong/sssplitmaker-faq.html"
							if err := exec.CommandContext(context.Background(), "rundll32", "url.dll,FileProtocolHandler", url).Start(); err != nil {
								walk.MsgBox(mainWindow, "错误", "打开浏览器失败，请手动访问："+url, walk.MsgBoxIconError)
							}
						},
					},
					PushButton{Text: "清空", OnClicked: func() {
						resetLines(1)
					}},
					PushButton{Text: "另存为", OnClicked: onSaveSplitsFile},
				},
			},
		},
	}.Create()
	addLine()
	if err != nil {
		walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
		return
	}
	hWnd := mainWindow.Handle()
	currStyle := win.GetWindowLong(hWnd, win.GWL_STYLE)
	win.SetWindowLong(hWnd, win.GWL_STYLE, currStyle & ^win.WS_SIZEBOX)
	mainWindow.Run()
}
