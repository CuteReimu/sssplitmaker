# 丝之歌计时器生成器（中文版）

![](https://img.shields.io/github/go-mod/go-version/CuteReimu/sssplitmaker "Language")
[![](https://img.shields.io/github/actions/workflow/status/CuteReimu/sssplitmaker/golangci-lint.yml?branch=master)](https://github.com/CuteReimu/sssplitmaker/actions/workflows/golangci-lint.yml "Analysis")
[![](https://img.shields.io/github/license/CuteReimu/sssplitmaker)](https://github.com/CuteReimu/sssplitmaker/blob/master/LICENSE "LICENSE")

## 如何使用

https://cutereimu.cn/daily/silksong/sssplitmaker-faq.html

## 编译说明

首先需要 Go 和 Nodejs，然后安装 wails：

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

然后使用以下命令就可以调试或打包了：

```bash
# 提前下载wasm文件
curl -O https://github.com/AlexKnauth/silksong-autosplit-wasm/releases/latest/download/silksong_autosplit_wasm_stable.wasm
curl -O https://raw.githubusercontent.com/LiveSplit/LiveSplit.AutoSplitters/master/LiveSplit.AutoSplitters.xml

# 本地调试
wails dev

# 打包
wails build -platform=windows/amd64 -webview2 embed
```

## （开发向）如何更新最新的节点

```bash
cd translate
./update_splits.sh
go test
```

执行完毕后会输出缺少了什么节点，将其补全。

打开 `main.go` 搜索 `Auto Splitter Version`，更新一下版本号。之后执行`./build.sh`重新编译即可。

## 特别鸣谢

感谢 AlexKnauth 大佬编写的丝之歌 AutoSplitter： https://github.com/AlexKnauth/silksong-autosplit-wasm

默认模板和图标来自： https://github.com/slaurent22/hk-split-maker
