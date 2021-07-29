package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/tools/txtar"
)

func TestStartWork(t *testing.T) {
	fis, err := ioutil.ReadDir("testdata")
	if err != nil {
		t.Fatal(err)
	}
	for _, fi := range fis {
		name := strings.TrimSuffix(filepath.Base(fi.Name()), ".txt")
		t.Run(name, func(t *testing.T) { // definitely can't run these in parallel...
			dir := t.TempDir()
			if err != nil {
				t.Fatal(err)
			}
			ar, err := txtar.ParseFile(filepath.Join("testdata", fi.Name()))
			if err != nil {
				t.Fatal(err)
			}
			var want string
			for _, f := range ar.Files {
				if f.Name == "go.work.want" {
					want = string(f.Data)
				}
				path := filepath.Join(dir, f.Name)
				if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, f.Data, 0666); err != nil {
					t.Fatal(err)
				}
			}
			if want == "" {
				t.Fatal("failed to find go.work.want file")
			}
			oldwd, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}
			defer os.Chdir(oldwd)
			os.Chdir(dir)

			startWork()
			gotb, err := ioutil.ReadFile("go.work")
			if err != nil {
				t.Fatalf("reading go.work output file from startwork: %s", err)
			}
			got := string(gotb)
			if got != want {
				t.Errorf("go.work output: got %s, want %s", got, want)
			}
		})
	}
}
