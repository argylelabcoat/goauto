// Copyright 2021 Matthew Hughes. All rights reserved.
// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

// Package gotask implements tasks for running Go specific tools
package gotask

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/argylelabcoat/goauto"
)

type goPrjTask struct {
	gocmd string
	args  []string
}

// NewGoPrjTask returns a goauto.Tasker that runs a go command with optional arguments
// All of the commands are run on the project directory associated with goauto.TaskInfo.Target
// goauto.TaskInfo.Target is not updated
func NewGoPrjTask(gocmd string, args ...string) goauto.Tasker {
	return &goPrjTask{gocmd: gocmd, args: args}
}

func (gt *goPrjTask) Run(info *goauto.TaskInfo) (err error) {
	t0 := time.Now()
	info.Target = info.Src
	info.Buf.Reset()
	dir := goauto.GoRelSrcDir(info.Src)
	targs := append([]string{gt.gocmd}, gt.args...)
	targs = append(targs, dir)
	gocmd := exec.Command("go", targs...)
	gocmd.Stdout = &info.Buf
	gocmd.Stderr = info.Terr
	defer func() {
		fmt.Fprint(info.Tout, info.Buf.String())
		if err == nil && info.Verbose {
			t1 := time.Now()
			fmt.Fprintf(info.Tout, "<< Go %v %v %v\n", strings.Title(gt.gocmd), dir, t1.Sub(t0))
		}
	}()
	return gocmd.Run()
}

// NewGoTestTask returns a new task that will run all the project tests
func NewGoTestTask(args ...string) goauto.Tasker {
	return &goPrjTask{gocmd: "test", args: args}
}

// NewGoVetTask returns a new task that will vet the project
func NewGoVetTask(args ...string) goauto.Tasker {
	return &goPrjTask{gocmd: "vet", args: args}
}

// NewGoBuildTask returns a task that will build the project
func NewGoBuildTask(args ...string) goauto.Tasker {
	return &goPrjTask{gocmd: "build", args: args}
}

// NewGoInstallTask returns a task that will install the project
func NewGoInstallTask(args ...string) goauto.Tasker {
	return &goPrjTask{gocmd: "install", args: args}
}
