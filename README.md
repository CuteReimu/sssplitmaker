# 丝之歌计时器生成器（中文版）

![](https://img.shields.io/github/go-mod/go-version/CuteReimu/sssplitmaker "Language")
[![](https://img.shields.io/github/actions/workflow/status/CuteReimu/sssplitmaker/golangci-lint.yml?branch=master)](https://github.com/CuteReimu/sssplitmaker/actions/workflows/golangci-lint.yml "Analysis")
[![](https://img.shields.io/github/license/CuteReimu/sssplitmaker)](https://github.com/CuteReimu/sssplitmaker/blob/master/LICENSE "LICENSE")

> [!Note]
> 本项目只能在Windows环境下运行。

## 如何使用

https://cutereimu.cn/daily/silksong/sssplitmaker-faq.html

## 编译说明

**根据自己的编译环境，运行`build.bat`或`build.sh`即可进行编译。**

如果想要自己使用`go build`进行编译，需要提前下载wasm文件：

```shell
curl -O https://github.com/AlexKnauth/silksong-autosplit-wasm/releases/latest/download/silksong_autosplit_wasm_stable.wasm
curl -O https://raw.githubusercontent.com/LiveSplit/LiveSplit.AutoSplitters/master/LiveSplit.AutoSplitters.xml

# -ldflags中，-s是去掉符号表，-w是去掉调试信息，均可减小所生成二进制文件的体积
# -H=windowsgui是打开Windows窗口时隐藏控制台的黑框框
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -H=windowsgui" -o sssplitmaker.exe
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

默认模板来自： https://github.com/slaurent22/hk-split-maker
