package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	path := flag.String("path", "./", "replace path")
	ff := flag.String("file", "", "file format")
	from := flag.String("from", "", "old string")
	to := flag.String("to", "", "new string")
	thread := flag.Int("thread", 32, "replace thread")
	v := flag.Bool("v", false, "show info")
	showVersion := flag.Bool("version", false, "show version")

	flag.Parse()

	if *showVersion {
		fmt.Println(Version)
		return
	}

	if *to == "" || *from == "" {
		flag.Usage()
		return
	}

	if *thread < 1 {
		*thread = 1
	}

	sem := make(chan struct{}, *thread)
	var wg sync.WaitGroup

	err := filepath.Walk(*path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info == nil || info.IsDir() {
			return nil
		}

		if *ff != "" && !strings.HasSuffix(info.Name(), *ff) {
			return nil
		}

		wg.Add(1)
		sem <- struct{}{}
		go func(path string) {
			defer wg.Done()
			defer func() { <-sem }()
			replace(path, *from, *to, *v)
		}(path)

		return nil
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	wg.Wait()
}

func replace(f string, from string, to string, show bool) {
	data, err := os.ReadFile(f)
	if err != nil {
		panic(err)
	}

	str := strings.ReplaceAll(string(data), from, to)

	out, err := os.Create(f)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	if _, err := out.WriteString(str); err != nil {
		panic(err)
	}

	if show {
		fmt.Println("done " + f)
	}
}
