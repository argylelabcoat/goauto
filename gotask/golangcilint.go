// Copyright 2021 Matthew Hughes. All rights reserved.
// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

// Package gotask implements tasks for running Go specific tools
package gotask

import (
	"errors"
	"fmt"
	"os/exec"
	"time"

	"github.com/argylelabcoat/goauto"
)

type GolangciLintTask struct {
	args []string
}

func (lt *GolangciLintTask) Run(info *goauto.TaskInfo) (err error) {
	t0 := time.Now()
	info.Target = info.Src
	info.Buf.Reset()
	dir := goauto.GoRelSrcDir(info.Src)
	targs := append([]string{"run"}, lt.args...)
	cmd := exec.Command("golangci-lint", targs...)
	cmd.Stdout = &info.Buf
	cmd.Stderr = info.Terr
	defer func() {
		fmt.Fprint(info.Tout, info.Buf.String())
		if err == nil && info.Verbose {
			t1 := time.Now()
			fmt.Fprintf(info.Tout, ">>> golangci-lint %v %v\n", dir, t1.Sub(t0))
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
func NewGolangciLintTask(args ...string) goauto.Tasker {
	return &GolangciLintTask{args: args}
}
