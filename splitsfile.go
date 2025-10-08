package main

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"

	"github.com/CuteReimu/sssplitmaker/splitmaker"
	"github.com/CuteReimu/sssplitmaker/translate"
	"github.com/lxn/walk"
	"github.com/lxn/win"
)

type autoSplittingRuntimeSettings struct {
	Version        string
	CustomSettings []*xmlWasmSetting `xml:"CustomSettings>Setting"`
}

type xmlWasmSetting struct {
	Id      string `xml:"id,attr,omitempty"`
	Type    string `xml:"type,attr"`
	Value   string `xml:"value,attr,omitempty"`
	Setting []*xmlWasmSetting
}

type xmlRun struct {
	XMLName              xml.Name      `xml:"Run"`
	Version              string        `xml:"version,attr"`
	GameIcon             string        `xml:"GameIcon"`
	GameName             string        `xml:"GameName"`
	CategoryName         string        `xml:"CategoryName"`
	Offset               string        `xml:"Offset"`
	AttemptCount         int           `xml:"AttemptCount"`
	Segments             []*xmlSegment `xml:"Segments>Segment"`
	AutoSplitterSettings autoSplittingRuntimeSettings
	Other                []*xmlElement `xml:",any"`
}

type xmlSegment struct {
	Name  string
	Other []*xmlElement `xml:",any"`
}

type xmlElement struct {
	XMLName xml.Name
	Content string `xml:",innerxml"`
}

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

func loadLayoutFileFromSplitmaker(fileName string) {
	splitIds, err := splitmaker.GetSplitIds(fileName)
	if err != nil {
		walk.MsgBox(mainWindow, "获取splitmaker失败", err.Error(), walk.MsgBoxIconError)
		return
	}
	resetLines(len(splitIds) - 1)
	for i, id := range splitIds {
		index := translate.GetIndexByID(id)
		if index < 0 {
			walk.MsgBox(mainWindow, "解析Settings失败", "无法识别的分割点ID："+id, walk.MsgBoxIconError)
			return
		}
		if i == 0 {
			err = startTriggerComboBox.SetCurrentIndex(index)
		} else {
			name := translate.GetSplitDescriptionByID(id)
			splitIndex := strings.LastIndex(name, "（")
			if splitIndex > 0 {
				name = name[:splitIndex]
			}
			err = lines[i-1].name.SetText(name)
			err = lines[i-1].splitId.SetCurrentIndex(index)
		}
		if err != nil {
			walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
			return
		}
	}
}

func onSaveSplitsFile() {
	if fileRunData.Version == "" {
		fileRunData.Version = "1.7.0"
	}
	if fileRunData.AutoSplitterSettings.Version == "" {
		fileRunData.AutoSplitterSettings.Version = "1.0"
	}
	if fileRunData.GameName == "" {
		fileRunData.GameName = "Hollow Knight: Silksong"
	}
	if fileRunData.Offset == "" {
		fileRunData.Offset = "00:00:00"
	}
	splits := &xmlWasmSetting{
		Id:   "splits",
		Type: "list",
		Setting: []*xmlWasmSetting{{
			Type:  "string",
			Value: translate.GetIDByDescription(startTriggerComboBox.Text()),
		}},
	}
	fileRunData.AutoSplitterSettings.CustomSettings = []*xmlWasmSetting{{
		Id:    "script_name",
		Type:  "string",
		Value: "silksong_autosplit_wasm",
	}, splits}
	fileRunData.Segments = nil
	for _, line := range lines {
		splits.Setting = append(splits.Setting, &xmlWasmSetting{
			Type:  "string",
			Value: translate.GetIDByDescription(line.splitId.Text()),
		})
		fileRunData.Segments = append(fileRunData.Segments, &xmlSegment{
			Name:  line.name.Text(),
			Other: line.xmlSegmentOther,
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
