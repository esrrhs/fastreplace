package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestReplace(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "sample.txt")
	if err := os.WriteFile(file, []byte("aaa bbb aaabbb"), 0o644); err != nil {
		t.Fatal(err)
	}

	replace(file, "aaa", "xxx", false)

	got, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	want := "xxx bbb xxxbbb"
	if string(got) != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestReplaceRespectsSuffixFilterViaWalk(t *testing.T) {
	dir := t.TempDir()
	txt := filepath.Join(dir, "a.txt")
	goFile := filepath.Join(dir, "a.go")
	if err := os.WriteFile(txt, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(goFile, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}

	replace(txt, "old", "new", false)

	gotTxt, err := os.ReadFile(txt)
	if err != nil {
		t.Fatal(err)
	}
	gotGo, err := os.ReadFile(goFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(gotTxt) != "new" {
		t.Fatalf("txt = %q, want new", gotTxt)
	}
	if string(gotGo) != "old" {
		t.Fatalf("go file should be untouched, got %q", gotGo)
	}
}

func TestVersionNotEmpty(t *testing.T) {
	if Version == "" {
		t.Fatal("Version should not be empty")
	}
}

func TestRunReplacesWithoutDryRun(t *testing.T) {
	dir := t.TempDir()
	txt := filepath.Join(dir, "a.txt")
	goFile := filepath.Join(dir, "a.go")
	nested := filepath.Join(dir, "sub", "b.txt")
	if err := os.Mkdir(filepath.Join(dir, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(txt, []byte("old old"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(goFile, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(nested, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if err := run(&buf, dir, ".txt", "old", "new", 4, false, false); err != nil {
		t.Fatal(err)
	}
	if buf.Len() != 0 {
		t.Fatalf("unexpected output without -v: %q", buf.String())
	}

	gotTxt, err := os.ReadFile(txt)
	if err != nil {
		t.Fatal(err)
	}
	gotGo, err := os.ReadFile(goFile)
	if err != nil {
		t.Fatal(err)
	}
	gotNested, err := os.ReadFile(nested)
	if err != nil {
		t.Fatal(err)
	}
	if string(gotTxt) != "new new" {
		t.Fatalf("txt = %q, want %q", gotTxt, "new new")
	}
	if string(gotGo) != "old" {
		t.Fatalf("go file should be untouched, got %q", gotGo)
	}
	if string(gotNested) != "new" {
		t.Fatalf("nested txt = %q, want %q", gotNested, "new")
	}
}

func TestDryRunHits(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "sample.txt")
	content := []byte("aaa bbb aaabbb")
	before := writeAndStat(t, file, content)

	var buf bytes.Buffer
	if err := run(&buf, dir, "", "aaa", "xxx", 2, false, true); err != nil {
		t.Fatal(err)
	}

	assertFileUnchanged(t, file, before, content)

	want := fmt.Sprintf("%s: 2\n1 files, 2 replacements\n", file)
	if buf.String() != want {
		t.Fatalf("output = %q, want %q", buf.String(), want)
	}
}

func TestDryRunNoHits(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "sample.txt")
	content := []byte("aaa bbb aaabbb")
	before := writeAndStat(t, file, content)

	var buf bytes.Buffer
	if err := run(&buf, dir, "", "zzz", "xxx", 2, false, true); err != nil {
		t.Fatal(err)
	}

	assertFileUnchanged(t, file, before, content)

	want := "0 files, 0 replacements\n"
	if buf.String() != want {
		t.Fatalf("output = %q, want %q", buf.String(), want)
	}
}

func TestDryRunMultiFileAndSuffix(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.txt")
	b := filepath.Join(dir, "sub", "b.txt")
	c := filepath.Join(dir, "c.go")
	d := filepath.Join(dir, "d.txt")
	if err := os.Mkdir(filepath.Join(dir, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	aContent := []byte("old old")
	bContent := []byte("keep old keep")
	cContent := []byte("old old old")
	dContent := []byte("no match here")
	aInfo := writeAndStat(t, a, aContent)
	bInfo := writeAndStat(t, b, bContent)
	cInfo := writeAndStat(t, c, cContent)
	dInfo := writeAndStat(t, d, dContent)

	var buf bytes.Buffer
	if err := run(&buf, dir, ".txt", "old", "new", 4, false, true); err != nil {
		t.Fatal(err)
	}

	assertFileUnchanged(t, a, aInfo, aContent)
	assertFileUnchanged(t, b, bInfo, bContent)
	assertFileUnchanged(t, c, cInfo, cContent)
	assertFileUnchanged(t, d, dInfo, dContent)

	want := fmt.Sprintf("%s: 2\n%s: 1\n2 files, 3 replacements\n", a, b)
	if buf.String() != want {
		t.Fatalf("output = %q, want %q", buf.String(), want)
	}
}

func TestDryRunNonOverlappingCount(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "overlap.txt")
	content := []byte("aaa")
	before := writeAndStat(t, file, content)

	var buf bytes.Buffer
	if err := run(&buf, dir, "", "aa", "x", 1, false, true); err != nil {
		t.Fatal(err)
	}

	assertFileUnchanged(t, file, before, content)

	want := fmt.Sprintf("%s: 1\n1 files, 1 replacements\n", file)
	if buf.String() != want {
		t.Fatalf("output = %q, want %q", buf.String(), want)
	}
}

func writeAndStat(t *testing.T, path string, content []byte) os.FileInfo {
	t.Helper()
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	return info
}

func assertFileUnchanged(t *testing.T, path string, before os.FileInfo, content []byte) {
	t.Helper()
	after, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(content) {
		t.Fatalf("%s content changed: got %q, want %q", path, got, content)
	}
	if after.Size() != before.Size() {
		t.Fatalf("%s size changed: got %d, want %d", path, after.Size(), before.Size())
	}
	if after.Mode() != before.Mode() {
		t.Fatalf("%s mode changed: got %v, want %v", path, after.Mode(), before.Mode())
	}
	if !after.ModTime().Equal(before.ModTime()) {
		t.Fatalf("%s mtime changed: got %v, want %v", path, after.ModTime(), before.ModTime())
	}
}
