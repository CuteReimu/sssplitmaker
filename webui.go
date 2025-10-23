package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"regexp"
	"strings"
	"text/template"

	"github.com/CuteReimu/sssplitmaker/translate"
)

var (
	//go:embed static
	htmlFiles embed.FS
	//go:embed index.html
	htmlIndex string
)

func initWebUi() {
	t, err := template.New("result").Parse(htmlIndex)
	if err != nil {
		panic(t)
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

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b.Bytes())
	})
	mux.HandleFunc("/get-templates", func(w http.ResponseWriter, r *http.Request) {
		files := GetAllFiles()
		ret := make([]Option, 0, len(files))
		for _, f := range files {
			ret = append(ret, Option{Value: f, Label: f})
		}
		buf, err := json.Marshal(ret)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"code": -1, "msg": "internal server error"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(buf)
	})
	mux.HandleFunc("/get-splits", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"code": -1, "msg": "invalid params"}`))
			return
		}
		name := r.Form.Get("name")
		name, splits, err := GetSplitIds(name)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"code": -2, "msg": "read template failed"}`))
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
		buf, err := json.Marshal(map[string]any{
			"name":   name,
			"splits": retSplits,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"code": -1, "msg": "internal server error"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(buf)
	})
	mux.HandleFunc("/build-splits", webBuildSplits)
	mux.Handle("/x/", http.StripPrefix("/x/", http.FileServer(http.FS(htmlFiles))))

	go func() {
		if err := http.ListenAndServe("127.0.0.1:12333", mux); err != nil { //nolint:gosec
			panic(err)
		}
	}()
}

type webSplitLine struct {
	Name  string `json:"name"`
	Event string `json:"event"`
}

func webBuildSplits(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"code": -1, "msg": "invalid params"}`))
		return
	}

	var lines []webSplitLine
	data := r.Form.Get("data")
	err := json.Unmarshal([]byte(data), &lines)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"code": -1, "msg": "unmarshal failed"}`))
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
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"code": -3, "msg": "internal server error"}`))
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=splits.lss")
	w.Header().Set("Content-Type", "text-xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf)
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
