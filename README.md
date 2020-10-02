# fastreplace
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
