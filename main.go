package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
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
	dryRun := flag.Bool("dry-run", false, "preview matches without writing files")

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

	if err := run(os.Stdout, *path, *ff, *from, *to, *thread, *v, *dryRun); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type fileResult struct {
	path  string
	count int
}

func run(out io.Writer, root, suffix, from, to string, thread int, verbose, dryRun bool) error {
	if thread < 1 {
		thread = 1
	}

	sem := make(chan struct{}, thread)
	var wg sync.WaitGroup

	var mu sync.Mutex
	var results []fileResult

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info == nil || info.IsDir() {
			return nil
		}

		if suffix != "" && !strings.HasSuffix(info.Name(), suffix) {
			return nil
		}

		wg.Add(1)
		sem <- struct{}{}
		go func(path string) {
			defer wg.Done()
			defer func() { <-sem }()
			if dryRun {
				n := countReplacements(path, from)
				if n > 0 {
					mu.Lock()
					results = append(results, fileResult{path: path, count: n})
					mu.Unlock()
				}
				return
			}
			replace(path, from, to, verbose)
		}(path)

		return nil
	})
	if err != nil {
		return err
	}

	wg.Wait()

	if dryRun {
		printDryRun(out, results)
	}
	return nil
}

func countReplacements(f string, from string) int {
	data, err := os.ReadFile(f)
	if err != nil {
		panic(err)
	}
	return strings.Count(string(data), from)
}

func printDryRun(out io.Writer, results []fileResult) {
	sort.Slice(results, func(i, j int) bool { return results[i].path < results[j].path })
	files := 0
	reps := 0
	for _, r := range results {
		fmt.Fprintf(out, "%s: %d\n", r.path, r.count)
		files++
		reps += r.count
	}
	fmt.Fprintf(out, "%d files, %d replacements\n", files, reps)
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
