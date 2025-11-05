//go:build windows

package main

import (
	"slices"
	"strings"

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
	line            *walk.Composite
	name            *walk.LineEdit
	splitId         *walk.ComboBox
	icon            string
	xmlSegmentOther []*xmlElement
}

var lines []*lineData

func clearLine(line *lineData) {
	if err := line.name.SetText(""); err != nil {
		walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
	}
	if err := line.splitId.SetCurrentIndex(0); err != nil {
		walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
	}
	line.xmlSegmentOther = nil
}

func buildLine(line *lineData, c *Composite) {
	c.AssignTo = &line.line
	c.Layout = HBox{}
	c.MaxSize = Size{Width: 0, Height: 25}
	c.Children = []Widget{
		LineEdit{AssignTo: &line.name, MinSize: Size{Width: 200}, ToolTipText: "片段名"},
		ComboBox{AssignTo: &line.splitId, MinSize: Size{Width: 200},
			Model: splitDescriptions, Value: splitDescriptions[0],
			OnCurrentIndexChanged: func() {
				name := line.name.Text()
				isSubSplit := strings.HasPrefix(name, "-")
				name = strings.TrimLeft(name, "-")
				if name == "" || slices.ContainsFunc(translate.SplitsCache, func(d *translate.SplitData) bool {
					return strings.Contains(d.Description, name)
				}) {
					splitId := line.splitId.Text()
					splitIndex := strings.LastIndex(splitId, "（")
					if splitIndex > 0 {
						splitId = splitId[:splitIndex]
					}
					if isSubSplit {
						splitId = "-" + splitId
					}
					if err := line.name.SetText(splitId); err != nil {
						walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
					}
				}
			},
		},
		PushButton{Text: "✘", MaxSize: Size{Width: 25}, ToolTipText: "删除", OnClicked: func() {
			if len(lines) > 1 {
				removeLine(line)
			}
		}},
		PushButton{Text: "↑+", MaxSize: Size{Width: 25}, ToolTipText: "在上方增加一行",
			OnClicked: func() {
				idx := splitLinesView.Children().Index(line.line)
				addLine()
				moveLine(idx)
				clearLine(lines[idx])
			},
		},
		PushButton{Text: "↓+", MaxSize: Size{Width: 25}, ToolTipText: "在下方增加一行",
			OnClicked: func() {
				idx := splitLinesView.Children().Index(line.line)
				addLine()
				moveLine(idx + 1)
				clearLine(lines[idx+1])
			},
		},
		PushButton{Text: "↑", MaxSize: Size{Width: 25}, ToolTipText: "上移一行",
			OnClicked: func() {
				idx := splitLinesView.Children().Index(line.line)
				swapLine(idx-1, idx)
			},
		},
		PushButton{Text: "↓", MaxSize: Size{Width: 25}, ToolTipText: "下移一行",
			OnClicked: func() {
				idx := splitLinesView.Children().Index(line.line)
				swapLine(idx, idx+1)
			},
		},
	}
}

func addLine() {
	line := new(lineData)
	c := Composite{}
	buildLine(line, &c)
	err := c.Create(NewBuilder(splitLinesView))
	if err != nil {
		walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
	}
	lines = append(lines, line)
}

func removeLine(line *lineData) {
	idx := splitLinesView.Children().Index(line.line)
	if idx < 0 {
		walk.MsgBox(mainWindow, "错误", "无法删除这一行", walk.MsgBoxIconError)
		return
	}
	err := splitLinesView.Children().RemoveAt(idx)
	if err != nil {
		walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
		return
	}
	line.line.Dispose()
	lines = append(lines[:idx], lines[idx+1:]...)
}

func resetLines(count int) {
	if count >= len(lines) && len(lines) >= count/3 {
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
		c := Composite{}
		buildLine(line, &c)
		composite.Children = append(composite.Children, c)
		lines = append(lines, line)
	}
	err = composite.Create(NewBuilder(splitLinesViewContainer))
	if err != nil {
		walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
		return
	}
}

func swapLine(index1, index2 int) {
	if index1 == index2 || index1 < 0 || index2 < 0 || index1 >= len(lines) || index2 >= len(lines) {
		return
	}
	name1 := lines[index1].name.Text()
	name2 := lines[index2].name.Text()
	splitId1 := lines[index1].splitId.CurrentIndex()
	splitId2 := lines[index2].splitId.CurrentIndex()
	if err := lines[index1].name.SetText(name2); err != nil {
		walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
		return
	}
	if err := lines[index2].name.SetText(name1); err != nil {
		walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
		return
	}
	if err := lines[index1].splitId.SetCurrentIndex(splitId2); err != nil {
		walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
		return
	}
	if err := lines[index2].splitId.SetCurrentIndex(splitId1); err != nil {
		walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
		return
	}
	lines[index1].xmlSegmentOther, lines[index2].xmlSegmentOther = lines[index2].xmlSegmentOther, lines[index1].xmlSegmentOther
}

func moveLine(index int) {
	var err error
	for i := len(lines) - 2; i >= index; i-- {
		err = lines[i+1].name.SetText(lines[i].name.Text())
		if err != nil {
			walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
			return
		}
		err = lines[i+1].splitId.SetCurrentIndex(lines[i].splitId.CurrentIndex())
		if err != nil {
			walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
			return
		}
		lines[i+1].xmlSegmentOther = lines[i].xmlSegmentOther
	}
}
