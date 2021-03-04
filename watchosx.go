// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.
// +build darwin

package goauto

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/fsnotify/fsevents"
)

type watchOSX struct {
	out         io.Writer
	eventStream *fsevents.EventStream
	done        chan struct{}
	send        chan ESlice
}

// NewWatchOSX returns a OSX specific file system watcher
func NewWatchOSX() Watcher {
	w := new(watchOSX)
	w.eventStream = &fsevents.EventStream{
		Paths: []string{},
		Flags: fsevents.FileEvents | fsevents.WatchRoot,
	}
	return w
}

func (w *watchOSX) SetVerbose(out io.Writer) {
	w.out = out
}

func (w *watchOSX) convertFlags(e fsevents.Event) Op {
	var f Op
	if e.Flags&fsevents.ItemCreated == fsevents.ItemCreated {
		f |= Create
	}
	if e.Flags&fsevents.ItemRemoved == fsevents.ItemRemoved {
		f |= Remove
	}
	if e.Flags&fsevents.ItemRenamed == fsevents.ItemRenamed {
		f |= Rename
	}
	if e.Flags&fsevents.ItemModified == fsevents.ItemModified {
		f |= Write
	}
	if e.Flags&fsevents.ItemInodeMetaMod == fsevents.ItemInodeMetaMod {
		f |= Chmod
	}
	return f
}

func (w *watchOSX) Start(latency time.Duration, paths []string) (<-chan ESlice, error) {
	w.done = make(chan struct{})
	c := make(chan ESlice)
	w.send = c
	w.eventStream.Paths = paths
	w.eventStream.Latency = latency

	if w.out != nil {
		for _, d := range paths {
			fmt.Fprintln(w.out, "Watching", d)
		}
	}

	w.eventStream.Start()
	go w.bufferEvents(c, latency)
	return c, nil
}

// bufferEvents watches for file events and batches them up based on a timer
// if the event distributer is busy it just keeps batching up events
// **Thanks to github.com/egonelbre for the suggestions and examples for batch events
func (w *watchOSX) bufferEvents(send chan<- ESlice, l time.Duration) {
	defer close(send)

	tick := time.Tick(l)
	buf := make(ESlice, 0, 10)
	var out chan<- ESlice

	for {
		select {
		// buffer the events
		case msg := <-w.eventStream.Events:
			for _, e := range msg {
				buf = append(buf, &Event{Path: e.Path, Op: w.convertFlags(e)})
				if w.out != nil {
					fmt.Fprintln(w.out, Event{Path: e.Path, Op: w.convertFlags(e)})
				}
			}
		// check if we have any events
		case <-tick:
			if len(buf) > 0 {
				out = send
			}
		// if nil skip, otherwise send when it's ready
		case out <- buf:
			buf = make(ESlice, 0, 10)
			out = nil
		case <-w.done:
			return
		}
	}
}

func (w *watchOSX) Stop() error {
	if w.done == nil || w.eventStream == nil || w.send == nil {
		return errors.New("Watcher not started or already stopped")
	}
	if w.out != nil {
		fmt.Fprintln(w.out, "Watcher stopped")
	}
	select {
	case <-w.done:
	default:
		close(w.done)
	}
	w.eventStream.Stop()
	return nil
}

func (w *watchOSX) Add(path string) error {
	w.eventStream.Paths = append(w.eventStream.Paths, path)
	return nil
}

func (w *watchOSX) Remove(path string) error {
	return nil
}
