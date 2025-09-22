package main

import (
	"github.com/CuteReimu/sssplitmaker/translate"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var splitDescriptions []string

func init() {
	for _, split := range translate.SplitsCache {
		splitDescriptions = append(splitDescriptions, split.Description)
	}
}

type lineData struct {
	line    *walk.Composite
	name    *walk.LineEdit
	splitId *walk.ComboBox
}

var lines []*lineData

func addLine() {
	line := new(lineData)
	c := Composite{
		AssignTo: &line.line,
		Layout:   HBox{},
		MaxSize:  Size{Width: 0, Height: 25},
		Children: []Widget{
			LineEdit{AssignTo: &line.name, MinSize: Size{Width: 200}, ToolTipText: "片段名", Enabled: false},
			ComboBox{AssignTo: &line.splitId, MinSize: Size{Width: 200},
				Model: splitDescriptions, Value: splitDescriptions[0],
			},
			PushButton{Text: "+", MaxSize: Size{Width: 25}, ToolTipText: "在此位置增加一行",
				OnClicked: func() {
					idx := splitLinesView.Children().Index(line.line)
					moveLine(idx, true)
				},
			},
			PushButton{Text: "-", MaxSize: Size{Width: 25}, ToolTipText: "删掉此行",
				OnClicked: func() {
					idx := splitLinesView.Children().Index(line.line)
					moveLine(idx, false)
				},
			},
		},
	}
	err := c.Create(NewBuilder(splitLinesView))
	if err != nil {
		walk.MsgBox(nil, "错误", err.Error(), walk.MsgBoxIconError)
	}
	lines = append(lines, line)
}

func resetLines(count int) {
	if count >= len(lines) {
		for i := len(lines); i < count; i++ {
			addLine()
		}
		return
	}
	err := splitLinesViewContainer.Children().RemoveAt(0)
	if err != nil {
		walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
		return
	}
	splitLinesView.Dispose()
	composite := Composite{
		AssignTo:  &splitLinesView,
		Alignment: AlignHCenterVNear,
		Layout:    VBox{},
	}
	lines = []*lineData{}
	for range count {
		line := new(lineData)
		composite.Children = append(composite.Children, Composite{
			AssignTo: &line.line,
			Layout:   HBox{},
			MaxSize:  Size{Width: 0, Height: 25},
			Children: []Widget{
				LineEdit{AssignTo: &line.name, MinSize: Size{Width: 200}, Enabled: false},
				ComboBox{AssignTo: &line.splitId, MinSize: Size{Width: 200},
					Model: splitDescriptions, Value: splitDescriptions[0],
				},
				PushButton{Text: "+", MaxSize: Size{Width: 25}, ToolTipText: "在此位置增加一行",
					OnClicked: func() {
						idx := splitLinesView.Children().Index(line.line)
						moveLine(idx, true)
					},
				},
				PushButton{Text: "-", MaxSize: Size{Width: 25}, ToolTipText: "删掉此行",
					OnClicked: func() {
						idx := splitLinesView.Children().Index(line.line)
						moveLine(idx, false)
					},
				},
			},
		})
		lines = append(lines, line)
	}
	err = composite.Create(NewBuilder(splitLinesViewContainer))
	if err != nil {
		walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
		panic(err)
	}
}

func moveLine(index int, down bool) {
	if !down {
		for i := index; i <= len(lines)-2; i++ {
			err := lines[i].splitId.SetCurrentIndex(lines[i+1].splitId.CurrentIndex())
			if err != nil {
				walk.MsgBox(nil, "错误", err.Error(), walk.MsgBoxIconError)
				return
			}
		}
		err := lines[len(lines)-1].splitId.SetCurrentIndex(translate.GetIndexByID("EndingSplit"))
		if err != nil {
			walk.MsgBox(nil, "错误", err.Error(), walk.MsgBoxIconError)
		}
		return
	}
	for i := len(lines) - 2; i >= index; i-- {
		err := lines[i+1].splitId.SetCurrentIndex(lines[i].splitId.CurrentIndex())
		if err != nil {
			walk.MsgBox(nil, "错误", err.Error(), walk.MsgBoxIconError)
			return
		}
	}
}
