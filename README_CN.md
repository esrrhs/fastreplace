# fastreplace

[<img src="https://img.shields.io/github/license/esrrhs/fastreplace">](https://github.com/esrrhs/fastreplace)
[<img src="https://img.shields.io/github/languages/top/esrrhs/fastreplace">](https://github.com/esrrhs/fastreplace)
[<img src="https://img.shields.io/github/v/release/esrrhs/fastreplace">](https://github.com/esrrhs/fastreplace/releases)
[<img src="https://img.shields.io/github/downloads/esrrhs/fastreplace/total">](https://github.com/esrrhs/fastreplace/releases)
[<img src="https://img.shields.io/docker/pulls/esrrhs/fastreplace">](https://hub.docker.com/repository/docker/esrrhs/fastreplace)
[<img src="https://img.shields.io/github/actions/workflow/status/esrrhs/fastreplace/go.yml?branch=master">](https://github.com/esrrhs/fastreplace/actions)

> **多线程文件内容批量替换命令行工具。**

[English](README.md) | [中文说明](README_CN.md)

---

## 简介

**fastreplace** 递归遍历目录，可按文件后缀过滤，把匹配到的字符串全部替换成新内容；内部用有界协程池并行处理，适合大批量文本替换。

## 安装

从 **[Releases](https://github.com/esrrhs/fastreplace/releases)** 下载预编译二进制，或自行编译：

```bash
git clone https://github.com/esrrhs/fastreplace.git
cd fastreplace
go build -o fastreplace .
```

Docker：

```bash
docker pull esrrhs/fastreplace
```

## 使用

遍历当前目录及子目录的所有 `.txt` 文件，把内容 `aaa` 改为 `bbb`：

```bash
./fastreplace -path ./ -file .txt -from aaa -to bbb
```

查看版本：

```bash
./fastreplace -version
```

更多参数：

```text
Usage of ./fastreplace:
  -file string
        file format（后缀过滤，如 .txt）
  -from string
        old string
  -path string
        replace path (default "./")
  -thread int
        replace thread (default 32)
  -to string
        new string
  -v    show info
  -version
        show version
```

## 打包发布

```bash
./pack.sh
```

跨平台压缩包输出到 `pack/`。在 `version.go` 中 bump 版本并推送到 `master` 后，CI 会自动创建 GitHub Release。

## 许可

MIT — 见 [LICENSE](LICENSE)。
