// -*- coding:utf-8; indent-tabs-mode:nil; -*-

// Copyright 2014, Wu Xi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Extended log pkg for go.
// Output log format:
//	 "LEVEL: DATE TIME: GOROUTINE_ID: FILE:LINE: LOG_CONTENT\n"
// eg:
//	 "NOTICE: 08-06 10:45:19.598: 12345: bvc.go:100: hello world"
package logex

import (
    "errors"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

type Level uint
const (
    NONE    Level = iota
    FATAL
    WARNING
    NOTICE
    TRACE
    DEBUG
    LEVEL_MAX
)
var levelStr = [...]string{"NONE", "FATAL", "WARNING", "NOTICE", "TRACE", "DEBUG"}

type Logger struct {
	mu  sync.Mutex
	out io.Writer
	buf []byte     // for accumulating text to write
	enabled [LEVEL_MAX]bool // enabled log level
}

func New(level Level, out io.Writer) *Logger {
    l := &Logger{out: out}
    for i := FATAL; i <= level && i < LEVEL_MAX; i++ {
        l.enabled[i] = true
    }
    return l
}

var logger = New(DEBUG, os.Stderr)

func (l *Logger) Output(level Level, calldepth int, s string) error {
    if level > DEBUG {
        return errors.New("wrong log level")
    } else if !l.enabled[level] {
        return nil
    }

    _, file, line, ok := runtime.Caller(calldepth)
    if !ok {
        file = "???"
        line = 0
    }
	now := time.Now()

	l.mu.Lock()
	defer l.mu.Unlock()

	l.buf = l.buf[:0]
	l.formatPrefix(level, now, file, line)
	l.buf = append(l.buf, s...)
	if len(s) > 0 && s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	_, err := l.out.Write(l.buf)
	return err
}

func (l *Logger) formatPrefix(level Level, t time.Time, file string, line int) {
    var buf *[]byte = &l.buf
    *buf = append(*buf, levelStr[level]...)
    *buf = append(*buf, ": "...)

    _, month, day := t.Date()
    itoa(buf, int(month), 2)
    *buf = append(*buf, '-')
    itoa(buf, day, 2)
    *buf = append(*buf, ' ')
    hour, min, sec := t.Clock()
    itoa(buf, hour, 2)
    *buf = append(*buf, ':')
    itoa(buf, min, 2)
    *buf = append(*buf, ':')
    itoa(buf, sec, 2)
    *buf = append(*buf, '.')
    itoa(buf, t.Nanosecond()/1e3, 3)
    *buf = append(*buf, ": "...)

    itoa(buf, int(goid()), -1)
    *buf = append(*buf, ": "...)

    short := file
    for i := len(file) - 1; i > 0; i-- {
        if file[i] == '/' {
            short = file[i+1:]
            break
        }
    }
    *buf = append(*buf, short...)
    *buf = append(*buf, ':')
    itoa(buf, line, -1)
    *buf = append(*buf, ": "...)
}

// Cheap integer to fixed-width decimal ASCII.  Give a negative width to avoid zero-padding.
// Knows the buffer has capacity.
func itoa(buf *[]byte, i int, wid int) {
	var u uint = uint(i)
	if u == 0 && wid <= 1 {
		*buf = append(*buf, '0')
		return
	}

	// Assemble decimal in reverse order.
	var b [32]byte
	bp := len(b)
	for ; u > 0 || wid > 0; u /= 10 {
		bp--
		wid--
		b[bp] = byte(u%10) + '0'
	}
	*buf = append(*buf, b[bp:]...)
}

// Goid returns the id of goroutine, defined in ./goid.c
func goid() int32
