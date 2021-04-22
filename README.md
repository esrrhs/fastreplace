# fastreplace

[<img src="https://img.shields.io/github/license/esrrhs/fastreplace">](https://github.com/esrrhs/fastreplace)
[<img src="https://img.shields.io/github/languages/top/esrrhs/fastreplace">](https://github.com/esrrhs/fastreplace)
[![Go Report Card](https://goreportcard.com/badge/github.com/esrrhs/fastreplace)](https://goreportcard.com/report/github.com/esrrhs/fastreplace)
[<img src="https://img.shields.io/github/v/release/esrrhs/fastreplace">](https://github.com/esrrhs/fastreplace/releases)
[<img src="https://img.shields.io/github/downloads/esrrhs/fastreplace/total">](https://github.com/esrrhs/fastreplace/releases)
[<img src="https://img.shields.io/docker/pulls/esrrhs/fastreplace">](https://hub.docker.com/repository/docker/esrrhs/fastreplace)
[<img src="https://img.shields.io/github/workflow/status/esrrhs/fastreplace/Go">](https://github.com/esrrhs/fastreplace/actions)

多线程文件替换工具

# 使用
* 遍历当前目录及子目录的所有txt文件，把内容aaa改为bbb
```
# ./fastreplace -path ./ -file .txt -from aaa -to bbb
```
* 更多参数参考-h
```
Usage of fastreplace.exe:
  -file string
        file format
  -from string
        old string
  -path string
        replace path (default "./")
  -thread int
        replace thread (default 32)
  -to string
        new string
  -v    show info
```
