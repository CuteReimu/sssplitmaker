package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"text/template"

	"github.com/CuteReimu/sssplitmaker/translate"
	"github.com/gin-gonic/gin"
)

var (
	//go:embed static
	htmlFiles embed.FS
	//go:embed index.html
	htmlIndex string
	//go:embed translate.html
	htmlTranslate string
)

func initWebUi() {
	t, err := template.New("result").Parse(htmlIndex)
	if err != nil {
		panic(err)
	}
	t2, err := template.New("translate").Parse(htmlTranslate)
	if err != nil {
		panic(err)
	}

	type Option struct {
		Value string `json:"value"`
		Label string `json:"label"`
	}

	options := make([]Option, 0, len(translate.SplitsCache))
	for _, v := range translate.SplitsCache {
		options = append(options, Option{Value: v.ID, Label: v.Description})
	}
	buf, err := json.MarshalIndent(options, "", "  ")
	if err != nil {
		panic(err)
	}
	var b bytes.Buffer
	err = t.Execute(&b, map[string]any{"options": string(buf)})
	if err != nil {
		panic(err)
	}

	var b2 bytes.Buffer
	err = t2.Execute(&b2, map[string]any{"tableData": translate.SplitsHtml})
	if err != nil {
		panic(err)
	}

	gin.SetMode(gin.ReleaseMode)
	g := gin.New()
	g.GET("/", func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", b.Bytes())
	})
	g.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": -3, "msg": fmt.Sprintf("获取文件失败: %+v", err)})
			return
		}
		f, err := file.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": -4, "msg": fmt.Sprintf("打开文件失败: %+v", err)})
			return
		}
		defer func() { _ = f.Close() }()
		buf, err := io.ReadAll(f)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": -5, "msg": fmt.Sprintf("读取文件失败: %+v", err)})
			return
		}
		result, err := webLoadSplitFile(buf)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": -2, "msg": fmt.Sprintf("解析文件失败: %+v", err)})
			return
		}
		c.JSON(http.StatusOK, result)
	})
	g.GET("/translate", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", b2.Bytes())
	})
	g.GET("/get-templates", func(c *gin.Context) {
		files := GetAllFiles()
		ret := make([]Option, 0, len(files))
		for _, f := range files {
			ret = append(ret, Option{Value: f, Label: f})
		}
		c.JSON(http.StatusOK, ret)
	})
	g.GET("/get-options", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": options})
	})
	g.GET("/get-splits", func(c *gin.Context) {
		name := c.Query("name")
		name, splits, err := GetSplitIds(name)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": -2, "msg": "read template failed"})
			return
		}
		retSplits := make([]webSplitLine, 0, len(splits)+1)
		for _, id := range splits {
			isSubSplit := strings.HasPrefix(id, "-")
			id = strings.TrimLeft(id, "-")
			id = regexp.MustCompile(`\{.*?}`).ReplaceAllString(id, "")
			name := translate.GetSplitDescriptionByID(id)
			splitIndex := strings.LastIndex(name, "（")
			if splitIndex > 0 {
				name = name[:splitIndex]
			}
			if isSubSplit {
				name = "-" + name
			}
			retSplits = append(retSplits, webSplitLine{Name: name, Event: id})
		}
		c.JSON(http.StatusOK, gin.H{"name": name, "splits": retSplits})
	})
	g.POST("/build-splits", webBuildSplits)
	g.StaticFS("/x/", http.FS(htmlFiles))

	go func() {
		if err = g.Run("127.0.0.1:12333"); err != nil {
			panic(err)
		}
	}()
}

type webSplitLine struct {
	Name  string `json:"name"`
	Event string `json:"event"`
}

func webBuildSplits(c *gin.Context) {
	var lines []webSplitLine
	data := c.PostForm("data")
	err := json.Unmarshal([]byte(data), &lines)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "unmarshal failed"})
		return
	}

	fileRunData := &xmlRun{
		Version:  "1.7.0",
		GameName: "Hollow Knight: Silksong",
		Offset:   "00:00:00",
		AutoSplitterSettings: autoSplittingRuntimeSettings{
			Version: "1.0",
		},
	}
	splits := &xmlWasmSetting{
		Id:   "splits",
		Type: "list",
	}
	fileRunData.AutoSplitterSettings.CustomSettings = []*xmlWasmSetting{{
		Id:    "script_name",
		Type:  "string",
		Value: "silksong_autosplit_wasm",
	}, splits}
	for i, line := range lines {
		splits.Setting = append(splits.Setting, &xmlWasmSetting{
			Type:  "string",
			Value: line.Event,
		})
		if i == 0 {
			continue
		}
		fileRunData.Segments = append(fileRunData.Segments, &xmlSegment{
			Name: line.Name,
		})
	}
	buf, err := xml.MarshalIndent(fileRunData, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -3, "msg": "internal server error"})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=splits.lss")
	c.Data(http.StatusOK, "text/xml; charset=utf-8", buf)
}

func webLoadSplitFile(buf []byte) ([]webSplitLine, error) {
	run := &xmlRun{}
	err := xml.Unmarshal(buf, run)
	if err != nil {
		return nil, err
	}

	result := []webSplitLine{{}}
	for _, segment := range run.Segments {
		result = append(result, webSplitLine{Name: segment.Name})
	}

	for _, setting := range run.AutoSplitterSettings.CustomSettings {
		if setting.Id == "splits" {
			if setting.Type != "list" {
				return nil, errors.New("splits字段类型错误")
			} else {
				for i, s := range setting.Setting {
					if s.Type != "string" {
						return nil, errors.New("splits子字段类型错误")
					} else {
						if i < len(result) {
							result[i].Event = s.Value
						} else {
							result = append(result, webSplitLine{Event: s.Value})
						}
					}
				}
			}
			break
		}
	}
	return append(result, webSplitLine{Name: "ManualSplit"}), nil
}

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
