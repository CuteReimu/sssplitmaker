package main

import (
	"context"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"os"
	"regexp"
	"strings"

	"github.com/CuteReimu/sssplitmaker/translate"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx         context.Context
	uploadedRun *xmlRun
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

type Option struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

func (a *App) GetOptions() []Option {
	options := make([]Option, 0, len(translate.SplitsCache))
	for _, v := range translate.SplitsCache {
		options = append(options, Option{Value: v.ID, Label: v.Description})
	}
	return options
}

func (a *App) GetTemplates() []Option {
	files := GetAllFiles()
	ret := make([]Option, 0, len(files))
	for _, f := range files {
		ret = append(ret, Option{Value: f, Label: f})
	}
	return ret
}

type SplitLine struct {
	Name  string        `json:"name"`
	Event string        `json:"event"`
	Icon  string        `json:"icon"`
	Other []*xmlElement `json:"other"`
}

type GetSplitsResult struct {
	Name   string      `json:"name"`
	Splits []SplitLine `json:"splits"`
}

// LoadSplitFile parses an .lss file from its XML text content
func (a *App) LoadSplitFile(content string) ([]SplitLine, error) {
	run := &xmlRun{}
	if err := xml.Unmarshal([]byte(content), run); err != nil {
		return nil, err
	}

	result := []SplitLine{{}}
	for _, segment := range run.Segments {
		result = append(result, SplitLine{
			Name:  segment.Name,
			Other: segment.Other,
			Icon:  convertIconToHtmlFormat(segment.Icon.Icon),
		})
	}

	for _, setting := range run.AutoSplitterSettings.CustomSettings {
		if setting.Id == "splits" {
			if setting.Type != "list" {
				return nil, errors.New("splits字段类型错误")
			}
			for i, s := range setting.Setting {
				if s.Type != "string" {
					return nil, errors.New("splits子字段类型错误")
				}
				if i < len(result) {
					result[i].Event = s.Value
				} else {
					result = append(result, SplitLine{Event: s.Value})
				}
			}
			break
		}
	}

	a.uploadedRun = run
	return append(result, SplitLine{Name: "ManualSplit"}), nil
}

// GetSplits returns split lines from a template file
func (a *App) GetSplits(name string) (*GetSplitsResult, error) {
	categoryName, ids, err := GetSplitIds(name)
	if err != nil {
		return nil, err
	}

	splits := make([]SplitLine, 0, len(ids))
	for i, id := range ids {
		isSubSplit := strings.HasPrefix(id, "-")
		id = strings.TrimLeft(id, "-")
		id = regexp.MustCompile(`\{.*?}`).ReplaceAllString(id, "")
		var splitName, icon string
		if i > 0 {
			splitName = translate.GetSplitDescriptionByID(id)
			if idx := strings.LastIndex(splitName, "（"); idx > 0 {
				splitName = splitName[:idx]
			}
			if isSubSplit {
				splitName = "-" + splitName
			}
			icon = getIconHtmlFormat(id)
		}
		splits = append(splits, SplitLine{Name: splitName, Event: id, Icon: icon})
	}

	return &GetSplitsResult{Name: categoryName, Splits: splits}, nil
}

// GetIcon returns the icon in HTML img-src format for a split ID
func (a *App) GetIcon(splitId string) string {
	return getIconHtmlFormat(splitId)
}

// buildSplits builds the LSS XML and returns base64-encoded bytes
func (a *App) buildSplits(data []SplitLine, includeTimeRecords bool) (string, error) {
	var fileRunData xmlRun
	if a.uploadedRun != nil {
		fileRunData = *a.uploadedRun
	}
	if !includeTimeRecords {
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
	splits := &xmlWasmSetting{Id: "splits", Type: "list"}
	fileRunData.AutoSplitterSettings.CustomSettings = []*xmlWasmSetting{
		{Id: "script_name", Type: "string", Value: "silksong_autosplit_wasm"},
		splits,
	}
	fileRunData.Segments = nil
	for i, line := range data {
		splits.Setting = append(splits.Setting, &xmlWasmSetting{Type: "string", Value: line.Event})
		if i == 0 {
			if splits.Setting[0].Value == "Act1Start" {
				fileRunData.Offset = "00:00:21.7600000"
			}
			continue
		}
		var other []*xmlElement
		if includeTimeRecords {
			other = line.Other
		}
		fileRunData.Segments = append(fileRunData.Segments, &xmlSegment{
			Name:  line.Name,
			Other: other,
			Icon:  xmlIcon{convertIconToLiveSplitFormat(line.Icon)},
		})
	}

	buf, err := xml.MarshalIndent(fileRunData, "", "  ")
	if err != nil {
		return "", err
	}
	buf = append([]byte(`<?xml version="1.0" encoding="UTF-8"?>`+"\n"), buf...)
	return base64.StdEncoding.EncodeToString(buf), nil
}

// SaveSplitsFile shows a native save dialog and writes the LSS file
func (a *App) SaveSplitsFile(data []SplitLine, includeTimeRecords bool) error {
	b64, err := a.buildSplits(data, includeTimeRecords)
	if err != nil {
		return err
	}
	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return err
	}
	dest, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		DefaultFilename: "splits.lss",
		Filters: []runtime.FileFilter{
			{DisplayName: "LiveSplit splits (*.lss)", Pattern: "*.lss"},
		},
	})
	if err != nil || dest == "" {
		return err
	}
	return os.WriteFile(dest, raw, 0644)
}

// SaveIconsZip shows a native save dialog and writes the icons zip
func (a *App) SaveIconsZip() error {
	b64, err := a.downloadIcons()
	if err != nil {
		return err
	}
	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return err
	}
	dest, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		DefaultFilename: "icons.zip",
		Filters: []runtime.FileFilter{
			{DisplayName: "ZIP archive (*.zip)", Pattern: "*.zip"},
		},
	})
	if err != nil || dest == "" {
		return err
	}
	return os.WriteFile(dest, raw, 0644)
}

// downloadIcons returns base64-encoded zip of all icons
func (a *App) downloadIcons() (string, error) {
	buf, err := zipIcons()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf), nil
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.setWindowSize()
}

func (a *App) setWindowSize() {
	// 1. 获取所有屏幕的信息
	screens, _ := runtime.ScreenGetAll(a.ctx)
	if len(screens) == 0 {
		return // 获取失败，使用默认配置
	}

	// 2. 找到主屏幕 (IsPrimary) 或当前屏幕 (IsCurrent)
	var targetScreen *runtime.Screen
	for _, screen := range screens {
		if screen.IsPrimary || screen.IsCurrent {
			targetScreen = &screen
			if screen.IsCurrent {
				break
			}
		}
	}
	// 如果没找到主屏幕，就回退使用第一个屏幕
	if targetScreen == nil && len(screens) > 0 {
		targetScreen = &screens[0]
	}
	if targetScreen == nil {
		return
	}

	// 3. 计算高度：屏幕总高度减去 100 像素
	newHeight := targetScreen.Size.Height - 100

	// 防止负高度（比如屏幕高度本身就小于 100，虽然很少见）
	if newHeight < 200 {
		newHeight = targetScreen.Size.Height
	}

	runtime.WindowSetSize(a.ctx, 1000, newHeight)
	runtime.WindowCenter(a.ctx)
}

// --- XML types ---

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
	Icon  xmlIcon
	Other []*xmlElement `xml:",any"`
}

type xmlElement struct {
	XMLName xml.Name
	Content string `xml:",innerxml"`
}

type xmlIcon struct {
	Icon string `xml:",cdata"`
}
