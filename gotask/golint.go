// Copyright 2021 Matthew Hughes. All rights reserved.
// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

// Package gotask implements tasks for running Go specific tools
package gotask

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"argylelabcoat/goauto"
)

type goLintTask struct {
	args []string
}

func (lt *goLintTask) Run(info *goauto.TaskInfo) (err error) {
	t0 := time.Now()
	info.Target = info.Src
	info.Buf.Reset()
	dir := goauto.GoRelSrcDir(info.Src)
	targs := append(lt.args, dir)
	cmd := exec.Command("golint", targs...)
	cmd.Stdout = &info.Buf
	cmd.Stderr = info.Terr
	defer func() {
		fmt.Fprint(info.Tout, info.Buf.String())
		if err == nil && info.Verbose {
			t1 := time.Now()
			fmt.Fprintf(info.Tout, ">>> Go Lint %v %v\n", dir, t1.Sub(t0))
		}
	}()
	if err = cmd.Run(); err != nil {
		return
	}
	if info.Buf.Len() > 0 {
		err = errors.New("FAIL")
		return
	}
	return
}

// NewGoLintTask returns a task that will golint the project
func NewGoLintTask(args ...string) goauto.Tasker {
	return &goLintTask{args: args}
}

type goMetaLinterTask struct {
	args []string
}

func (t *goMetaLinterTask) Run(info *goauto.TaskInfo) (err error) {
	t0 := time.Now()
	info.Target = info.Src
	info.Buf.Reset()
	dir := filepath.Dir(info.Src)
	targs := append(t.args, dir)
	cmd := exec.Command("gometalinter", targs...)
	cmd.Stdout = &info.Buf
	cmd.Stderr = info.Terr
	defer func() {
		fmt.Fprint(info.Tout, info.Buf.String())
		if err == nil && info.Verbose {
			t1 := time.Now()
			fmt.Fprintf(info.Tout, ">>> Go Meta Linter %v %v\n", dir, t1.Sub(t0))
		}
	}()
	if err = cmd.Run(); err != nil {
		println()
		return
	}
	if info.Buf.Len() > 0 {
		err = errors.New("FAIL")
		return
	}
	return
}

// NewGoMetaLinterTask returns a task that will run gometalinter for the project
// go get github.com/alecthomas/gometalinter
// "Concurrently run Go lint tools and normalise their output"
func NewGoMetaLinterTask(args ...string) goauto.Tasker {
	return &goMetaLinterTask{args: args}
}
