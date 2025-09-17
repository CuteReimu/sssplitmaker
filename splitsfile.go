package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lxn/walk"
	"github.com/lxn/win"
)

type xmlLayout struct {
	XMLName    xml.Name        `xml:"Layout"`
	Version    string          `xml:"version,attr"`
	Components []*xmlComponent `xml:"Components>Component"`
	Other      []*xmlElement   `xml:",any"`
}

type xmlElement struct {
	XMLName xml.Name
	Content string `xml:",innerxml"`
}

type xmlComponent struct {
	Path     string
	Settings *xmlTempSettings
	Other    []*xmlElement `xml:",any"`
}

type xmlTempSettings struct {
	Settings string `xml:",innerxml"`
}

type autoSplittingRuntimeSettings struct {
	XMLName        xml.Name `xml:"Settings"`
	Version        string
	ScriptPath     string
	CustomSettings []*xmlWasmSetting `xml:"CustomSettings>Setting"`
}

type xmlWasmSetting struct {
	Id      string `xml:"id,attr,omitempty"`
	Type    string `xml:"type,attr"`
	Value   string `xml:"value,attr,omitempty"`
	Setting []*xmlWasmSetting
}

type xmlRun struct {
	XMLName  xml.Name      `xml:"Run"`
	Version  string        `xml:"version,attr"`
	Segments []*xmlSegment `xml:"Segments>Segment"`
}

type xmlSegment struct {
	Name string
}

var fileLayoutPath string
var fileLayoutData *xmlLayout
var fileWasmSettings *autoSplittingRuntimeSettings
var fileWasmSettingsString *string

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

func onClickLoadLayoutFile() {
	dlg := new(walk.FileDialog)
	dlg.Title = "打开Layout文件"
	dlg.Filter = "Layout文件（*.lsl）|*.lsl"
	dlg.Flags = win.OFN_FILEMUSTEXIST
	if ok, err := dlg.ShowOpen(mainWindow); err != nil {
		walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
	} else if ok {
		file := dlg.FilePath
		if filepath.Ext(file) != ".lsl" {
			return
		}
		buf, err := os.ReadFile(file)
		if err != nil {
			walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
			return
		}
		loadLayoutFile(file, buf)
	}
}

func loadSplitFile(buf []byte) {
	run := &xmlRun{}
	err := xml.Unmarshal(buf, run)
	if err != nil {
		walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
		return
	}
	resetLines(max(len(lines), len(run.Segments)))
	for i, segments := range run.Segments {
		err = lines[i].name.SetText(segments.Name)
		if err != nil {
			walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
			return
		}
	}
}

func loadLayoutFile(file string, buf []byte) {
	run := &xmlLayout{}
	err := xml.Unmarshal(buf, run)
	if err != nil {
		walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
		return
	}
	for _, component := range run.Components {
		if component.Path == "LiveSplit.AutoSplittingRuntime.dll" {
			var settings autoSplittingRuntimeSettings
			err := xml.Unmarshal([]byte("<Settings>"+component.Settings.Settings+"</Settings>"), &settings)
			if err != nil {
				walk.MsgBox(mainWindow, "解析Settings失败", err.Error(), walk.MsgBoxIconError)
				return
			}
			for _, setting := range settings.CustomSettings {
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
								index := getIndexByID(s.Value)
								if index < 0 {
									walk.MsgBox(mainWindow, "解析Settings失败", fmt.Sprintf("无法识别的分割点ID：%s", s.Value), walk.MsgBoxIconError)
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
			fileWasmSettings = &settings
			fileWasmSettingsString = &component.Settings.Settings
			break
		}
	}
	fileLayoutPath = file
	fileLayoutData = run
	saveButton.SetEnabled(true)
}

func onSaveLayoutFile() {
	fileWasmSettings.CustomSettings = []*xmlWasmSetting{{
		Id:      "splits",
		Type:    "list",
		Setting: make([]*xmlWasmSetting, 0, len(lines)+1),
	}}
	fileWasmSettings.CustomSettings[0].Setting = append(fileWasmSettings.CustomSettings[0].Setting, &xmlWasmSetting{
		Type:  "string",
		Value: getIDByDescription(startTriggerComboBox.Text()),
	})
	for _, line := range lines {
		text := line.splitId.Text()
		id := getIDByDescription(text)
		fileWasmSettings.CustomSettings[0].Setting = append(fileWasmSettings.CustomSettings[0].Setting, &xmlWasmSetting{
			Type:  "string",
			Value: id,
		})
	}
	buf, err := xml.MarshalIndent(fileWasmSettings, "", "  ")
	if err != nil {
		walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
		return
	}
	fmt.Println(string(buf[len("<Settings>\n") : len(buf)-len("\n</Settings>")]))
	*fileWasmSettingsString = string(buf[len("<Settings>\n") : len(buf)-len("\n</Settings>")])
	if fileLayoutData == nil {
		walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
		return
	}
	buf, err = xml.MarshalIndent(fileLayoutData, "", "  ")
	if err != nil {
		walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
		return
	}
	dlg := new(walk.FileDialog)
	dlg.Title = "保存Layout文件"
	dlg.Filter = "Layout文件（*.lsl）|*.lsl"
	if ok, err := dlg.ShowSave(mainWindow); err != nil {
		walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
	} else if ok {
		file := dlg.FilePath
		if filepath.Ext(file) != ".lsl" {
			file += ".lsl"
		}
		err = os.WriteFile(file, buf, 0644)
		if err != nil {
			walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
			return
		}
	}
}
