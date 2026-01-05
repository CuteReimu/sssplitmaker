//go:build windows

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
var skipStartAnimationCheckBox *walk.CheckBox
var includeTimeRecordsCheckBox *walk.CheckBox

func main() {
	initWebUi()
	if walk.DlgCmdOK == walk.MsgBox(nil, "提示", "除“更新LiveSplit”功能以外，计时器生成器网页版在其它功能方面都更加好用，是否打开网页版？", walk.MsgBoxOKCancel|walk.MsgBoxIconInformation) {
		_ = exec.CommandContext(context.Background(), "rundll32", "url.dll,FileProtocolHandler", "http://127.0.0.1:12333/").Start()
	}
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
					PushButton{
						Text:      "打开网页版",
						Alignment: AlignHNearVCenter,
						OnClicked: func() {
							_ = exec.CommandContext(context.Background(), "rundll32", "url.dll,FileProtocolHandler", "http://127.0.0.1:12333/").Start()
						},
					},
					TextLabel{
						TextAlignment: AlignHFarVCenter,
						Text:          "Auto Splitter Version: 1.19.0",
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
						Text:          "由于翻译水平有限，翻译可能不准确，请检查翻译对照表是否正确后使用。",
					},
					PushButton{
						Text:      "打开翻译对照表",
						Alignment: AlignHFarVCenter,
						OnClicked: func() {
							_ = exec.CommandContext(context.Background(), "rundll32", "url.dll,FileProtocolHandler", "http://127.0.0.1:12333/translate").Start()
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
								AssignTo: &skipStartAnimationCheckBox,
								Text:     "跳过开局动画",
								OnClicked: func() {
									var err error
									id := translate.SplitsCache[startTriggerComboBox.CurrentIndex()].ID
									if id == "StartNewGame" || id == "Act1Start" {
										if skipStartAnimationCheckBox.Checked() {
											err = startTriggerComboBox.SetCurrentIndex(translate.GetIndexByID("Act1Start"))
										} else {
											err = startTriggerComboBox.SetCurrentIndex(translate.GetIndexByID("StartNewGame"))
										}
									}
									if err != nil {
										walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
									}
								},
							},
							ComboBox{
								AssignTo: &startTriggerComboBox,
								Model:    splitDescriptions,
								Value:    translate.GetSplitDescriptionByID("StartNewGame"),
								OnCurrentIndexChanged: func() {
									switch translate.SplitsCache[startTriggerComboBox.CurrentIndex()].ID {
									case "StartNewGame":
										skipStartAnimationCheckBox.SetEnabled(true)
										skipStartAnimationCheckBox.SetChecked(false)
									case "Act1Start":
										skipStartAnimationCheckBox.SetEnabled(true)
										skipStartAnimationCheckBox.SetChecked(true)
									default:
										skipStartAnimationCheckBox.SetEnabled(false)
									}
								},
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
			CheckBox{
				AssignTo:  &includeTimeRecordsCheckBox,
				Alignment: AlignHFarVCenter,
				Checked:   true,
				Text:      "保留*.lss文件中原本的时间记录（如果有）",
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
