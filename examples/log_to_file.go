// -*- coding:utf-8; indent-tabs-mode:nil; -*-

// Copyright 2014, Wu Xi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The very basic usage of logex pkg
package main

import (
	"fmt"
	"os"
)

import "github.com/wuxicn/logex"

func main() {
	if err := logex.SetUpFileLogger("./log", "example", nil); err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(255)
	}
	logex.SetLevel(logex.TRACE)

	fmt.Println("all logs will write to ./log/example.log file")

	logex.Debug("this message won't show")
	logex.Trace("this is trace message")
	logex.Notice("log levels: FATAL > WARNING > NOTICE > TRACE > DEBUG")
	logex.Warning("this is warning message")
}
