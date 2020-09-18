package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
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

	files, err := ioutil.ReadDir(*path)
	if err != nil {
		panic(err)
		return
	}

	var num int32
	for _, f := range files {

		if f == nil || f.IsDir() {
			continue
		}

		if *ff != "" {
			if !strings.HasSuffix(f.Name(), *ff) {
				continue
			}
		}

		if int(num) < *thread {
			atomic.AddInt32(&num, 1)
			go replace(f, &num, *from, *to, *v)
		} else {
			atomic.AddInt32(&num, 1)
			replace(f, &num, *from, *to, *v)
		}
	}

	for num > 0 {
		time.Sleep(time.Millisecond * 10)
	}
}

func replace(f os.FileInfo, num *int32, from string, to string, show bool) {

	defer atomic.AddInt32(num, -1)

	data, err := ioutil.ReadFile(f.Name())
	if err != nil {
		panic(err)
	}
	str := string(data)
	str = strings.Replace(str, from, to, -1)

	out, err := os.Create(f.Name())
	if err != nil {
		panic(err)
	}
	defer out.Close()

	out.WriteString(str)

	if show {
		fmt.Println("done " + f.Name())
	}
}
