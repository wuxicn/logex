// -*- coding:utf-8; indent-tabs-mode:nil; -*-

// Copyright 2014, Wu Xi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The very basic usage of logex pkg
package main

import "time"
import "github.com/wuxicn/logex"

func main() {
	logex.SetLevel(logex.NOTICE)
	logex.Debug("this message won't show")
	logex.Trace("this message won't show neither")
	logex.Notice("hi, pi is:", 3.1415926)
	logex.Warning("this is warning message")
	logex.Fatal("all logs will output to os.Stderr by default")

	d, _ := time.ParseDuration("10ms")

	c := make(chan bool)
	go func() {
		logex.Notice("note the third field of log line is the goroutine-id")
		time.Sleep(d)
		logex.Notice("this log is in another goroutine differs from main goroutine")
		c <- true
	}()

	f()

	time.Sleep(d)
	logex.Notice("in main")

	<-c
}

func f() {
	logex.Notice("log in f() function")
}
