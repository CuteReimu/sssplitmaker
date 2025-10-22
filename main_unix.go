//go:build unix

package main

import (
	"embed"
	_ "embed"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var (
	//go:embed static
	htmlFiles embed.FS
	//go:embed index.html
	htmlIndex string
)

func main() {
	initWebUi()

	time.Sleep(100 * time.Millisecond)

	switch strings.ToLower(runtime.GOOS) {
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
