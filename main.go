package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"
)

func main() {

	path := flag.String("path", "./", "replace path")
	ff := flag.String("file", "", "file format")
	from := flag.String("from", "", "old string")
	to := flag.String("to", "", "new string")
	thread := flag.Int("thread", 32, "replace thread")
	v := flag.Bool("v", false, "show info")

	flag.Parse()

	if *to == "" || *from == "" {
		flag.Usage()
		return
	}

	var num int32
	filepath.Walk(*path, func(path string, info os.FileInfo, err error) error {

		if info == nil || info.IsDir() {
			return nil
		}

		if *ff != "" {
			if !strings.HasSuffix(info.Name(), *ff) {
				return nil
			}
		}

		if int(num) < *thread {
			atomic.AddInt32(&num, 1)
			go replace(path, &num, *from, *to, *v)
		} else {
			atomic.AddInt32(&num, 1)
			replace(path, &num, *from, *to, *v)
		}

		return nil
	})

	for num > 0 {
		time.Sleep(time.Millisecond * 10)
	}
}

func replace(f string, num *int32, from string, to string, show bool) {

	defer atomic.AddInt32(num, -1)

	data, err := ioutil.ReadFile(f)
	if err != nil {
		panic(err)
	}
	str := string(data)
	str = strings.Replace(str, from, to, -1)

	out, err := os.Create(f)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	out.WriteString(str)

	if show {
		fmt.Println("done " + f)
	}
}
