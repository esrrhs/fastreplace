package main

import (
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
