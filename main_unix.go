//go:build unix

package main

import (
	"bytes"
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"text/template"
	"time"

	"github.com/CuteReimu/sssplitmaker/translate"
)

var (
	//go:embed static
	htmlFiles embed.FS
	//go:embed index.html
	htmlIndex string
)

func main() {
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
	buf, _ := json.MarshalIndent(options, "", "  ")
	var b bytes.Buffer
	err = t.Execute(&b, map[string]any{"options": string(buf)})
	if err != nil {
		fmt.Println(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(b.Bytes())
	})
	mux.Handle("/x/", http.StripPrefix("/x/", http.FileServer(http.FS(htmlFiles))))

	go func() {
		if err := http.ListenAndServe("127.0.0.1:12333", mux); err != nil {
			panic(err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	switch strings.ToLower(runtime.GOOS) {
	case "windows":
		_ = exec.Command("rundll32", "url.dll,FileProtocolHandler", "http://127.0.0.1:12333/").Start()
	case "linux":
		_ = exec.Command("xdg-open", "http://127.0.0.1:12333/").Start()
	case "darwin":
		_ = exec.Command("open", "http://127.0.0.1:12333/").Start()
	default:
		fmt.Println("不支持自动打开浏览器的操作系统：", runtime.GOOS)
		fmt.Println("请手动打开浏览器并访问: http://127.0.0.1:12333/")
	}

	select {}
}
