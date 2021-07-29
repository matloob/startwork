// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// startwork creates a go.work file containing all modules under
// the current working directory.
package main

import (
	"flag"
	"fmt"
	"go/build"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
)

var help = flag.Bool("help", false, "provides help instead of creating a go.work file")

func main() {

	if *help {
		fmt.Fprint(os.Stdout,
			`startwork creates a go.work file containing all the modules under
the current working directory. It expects there to not already be a go.work
file contained in the current working directory. (A future version may
support adding the modules under the current working directory to an
already existing go.work file.) It's intended to help easily set up a
go.work file for many modules, or to create a workspace similar to what
GOPATH mode provided.
`)
	}

	startWork()
}

func startWork() {
	// TODO(#45713) standardize paths

	workFile := "go.work"
	if _, err := os.Stat(workFile); !os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "startwork: go.work already exists in current directory.\nstartwork doesn't yet support editing already existing go.work files")
		os.Exit(1)
	}

	var modDirs []string
	filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if !info.Mode().IsDir() && filepath.Base(path) == "go.mod" {
			modDirs = append(modDirs, filepath.Dir(path))
		}

		return nil
	})

	goV := latestGoVersion() // Use current Go version by default
	workF := new(modfile.WorkFile)
	workF.Syntax = new(modfile.FileSyntax)
	workF.AddGoStmt(goV)

	for _, dir := range modDirs {
		// TODO(#45713): Add the module path of the module.
		workF.AddDirectory(dir, "")
	}

	data := modfile.Format(workF.Syntax)
	ioutil.WriteFile(workFile, data, 0644)
}

// latestGoVersion returns the latest version of the Go language supported by
// the toolchain that built this command.
func latestGoVersion() string {
	tags := build.Default.ReleaseTags
	version := tags[len(tags)-1]
	if !strings.HasPrefix(version, "go") || !modfile.GoVersionRE.MatchString(version[2:]) {
		fatalf("go: internal error: unrecognized default version %q", version)
	}
	return version[2:]
}

func fatalf(s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s, args...)
	os.Exit(1)
}
