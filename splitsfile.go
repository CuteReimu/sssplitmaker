//go:build windows

package main

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/CuteReimu/sssplitmaker/translate"
	"github.com/lxn/walk"
	"github.com/lxn/win"
)

var fileRunData = &xmlRun{}

func onClickLoadSplitFile() {
	dlg := new(walk.FileDialog)
	dlg.Title = "打开Splits文件"
	dlg.Filter = "Splits文件（*.lss）|*.lss"
	dlg.Flags = win.OFN_FILEMUSTEXIST
	if ok, err := dlg.ShowOpen(mainWindow); err != nil {
		walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
	} else if ok {
		file := dlg.FilePath
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
}

func loadSplitFile(buf []byte) {
	run := &xmlRun{}
	err := xml.Unmarshal(buf, run)
	if err != nil {
		walk.MsgBox(mainWindow, "解析文件失败", err.Error(), walk.MsgBoxIconError)
		return
	}
	if len(run.Segments) >= 50 {
		result := walk.MsgBox(mainWindow, "确认操作", "此lss文件的分段非常多，可能需要加载一段时间，是否继续？", walk.MsgBoxYesNo|walk.MsgBoxIconQuestion)
		if result != walk.DlgCmdYes {
			return
		}
	}

	resetLines(len(run.Segments))
	for i, segment := range run.Segments {
		err = lines[i].name.SetText(segment.Name)
		if err != nil {
			walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
			return
		}
		lines[i].xmlSegmentOther = segment.Other
	}

	for _, setting := range run.AutoSplitterSettings.CustomSettings {
		if setting.Id == "splits" {
			if setting.Type != "list" {
				walk.MsgBox(mainWindow, "解析Settings失败", "splits字段类型错误", walk.MsgBoxIconError)
				return
			} else {
				resetLines(max(len(lines), len(setting.Setting)-1))
				for i, s := range setting.Setting {
					if s.Type != "string" {
						walk.MsgBox(mainWindow, "解析Settings失败", "splits子字段类型错误", walk.MsgBoxIconError)
						return
					} else {
						index := translate.GetIndexByID(s.Value)
						if index < 0 {
							walk.MsgBox(mainWindow, "解析Settings失败", "无法识别的分割点ID："+s.Value, walk.MsgBoxIconError)
							return
						}
						if i == 0 {
							err = startTriggerComboBox.SetCurrentIndex(index)
						} else {
							err = lines[i-1].splitId.SetCurrentIndex(index)
						}
						if err != nil {
							walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
							return
						}
					}
				}
			}
			break
		}
	}
	fileRunData = run
}

var currentSplitmakerFileName string

func loadLayoutFileFromSplitmaker(fileName string) {
	if currentSplitmakerFileName == fileName {
		return
	}
	currentSplitmakerFileName = fileName
	categoryName, splitIds, err := GetSplitIds(fileName)
	if err != nil {
		walk.MsgBox(mainWindow, "获取splitmaker失败", err.Error(), walk.MsgBoxIconError)
		return
	}
	if len(splitIds) >= 50 {
		result := walk.MsgBox(mainWindow, "确认操作", "此模板的分段非常多，可能需要加载一段时间，是否继续？", walk.MsgBoxYesNo|walk.MsgBoxIconQuestion)
		if result != walk.DlgCmdYes {
			return
		}
	}
	resetLines(len(splitIds) - 1)
	for i, id := range splitIds {
		isSubSplit := strings.HasPrefix(id, "-")
		id = strings.TrimLeft(id, "-")
		id = regexp.MustCompile(`\{.*?}`).ReplaceAllString(id, "")
		index := translate.GetIndexByID(id)
		if index < 0 {
			walk.MsgBox(mainWindow, "解析Settings失败", "无法识别的分割点ID："+id, walk.MsgBoxIconError)
			return
		}
		if i == 0 {
			err = startTriggerComboBox.SetCurrentIndex(index)
			if err != nil {
				walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
				return
			}
		} else {
			name := translate.GetSplitDescriptionByID(id)
			splitIndex := strings.LastIndex(name, "（")
			if splitIndex > 0 {
				name = name[:splitIndex]
			}
			if isSubSplit {
				name = "-" + name
			}
			err = lines[i-1].name.SetText(name)
			if err != nil {
				walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
				return
			}
			err = lines[i-1].splitId.SetCurrentIndex(index)
			if err != nil {
				walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
				return
			}
		}
	}
	fileRunData = &xmlRun{CategoryName: categoryName}
}

func onSaveSplitsFile() {
	if !includeTimeRecordsCheckBox.Checked() {
		fileRunData.Other = nil
	}
	if fileRunData.Version == "" {
		fileRunData.Version = "1.7.0"
	}
	if fileRunData.AutoSplitterSettings.Version == "" {
		fileRunData.AutoSplitterSettings.Version = "1.0"
	}
	if fileRunData.GameName == "" {
		fileRunData.GameName = "Hollow Knight: Silksong"
	}
	fileRunData.Offset = "00:00:00"
	splits := &xmlWasmSetting{
		Id:   "splits",
		Type: "list",
		Setting: []*xmlWasmSetting{{
			Type:  "string",
			Value: translate.GetIDByDescription(startTriggerComboBox.Text()),
		}},
	}
	if splits.Setting[0].Value == "Act1Start" {
		fileRunData.Offset = "00:00:21.7600000"
	}
	fileRunData.AutoSplitterSettings.CustomSettings = []*xmlWasmSetting{{
		Id:    "script_name",
		Type:  "string",
		Value: "silksong_autosplit_wasm",
	}, splits}
	fileRunData.Segments = nil
	for _, line := range lines {
		splitId := translate.GetIDByDescription(line.splitId.Text())
		splits.Setting = append(splits.Setting, &xmlWasmSetting{
			Type:  "string",
			Value: splitId,
		})
		icon := line.icon
		if icon == "" {
			icon = getIcon(splitId)
		}
		var other []*xmlElement
		if includeTimeRecordsCheckBox.Checked() {
			other = line.xmlSegmentOther
		}
		fileRunData.Segments = append(fileRunData.Segments, &xmlSegment{
			Name:  line.name.Text(),
			Other: other,
			Icon:  xmlIcon{icon},
		})
	}
	buf, err := xml.MarshalIndent(fileRunData, "", "  ")
	if err != nil {
		walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
		return
	}
	buf = append([]byte(`<?xml version="1.0" encoding="UTF-8"?>`+"\n"), buf...)
	dlg := new(walk.FileDialog)
	dlg.Title = "保存分段文件"
	dlg.Filter = "分段文件（*.lss）|*.lss"
	dlg.Flags = win.OFN_OVERWRITEPROMPT | win.OFN_NOREADONLYRETURN
	if ok, err := dlg.ShowSave(mainWindow); err != nil {
		walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
	} else if ok {
		file := dlg.FilePath
		if filepath.Ext(file) != ".lss" {
			file += ".lss"
		}
		err = os.WriteFile(file, buf, 0644)
		if err != nil {
			walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
			return
		}
	}
}
