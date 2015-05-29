// -*- coding:utf-8; indent-tabs-mode:nil; -*-

// Copyright 2014, Wu Xi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Extended logging pkg for go.
// Output log format:
//	 "LEVEL: DATE TIME: g=GOROUTINE_ID: FILE:LINE: LOG_CONTENT\n"
// eg:
//	 "NOTICE: 08-06 10:45:19.598: g=12: bvc.go:100: hello world"
package logex

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

// Log Level type: FATAL > WARNING > NOTICE > TRACE > DEBUG
type Level uint

const (
	NONE Level = iota
	FATAL
	WARNING
	NOTICE
	TRACE
	DEBUG
	LEVEL_MAX
)

var levelStr = [...]string{"NONE", "FATAL", "WARNING", "NOTICE", "TRACE", "DEBUG"}

// A Logger represents an active logging object that generates lines of
// output to an io.Writer.  Each logging operation makes a single call to
// the Writer's Write method.  A Logger can be used simultaneously from
// multiple goroutines; it guarantees to serialize access to the Writer.
type Logger struct {
	mu      sync.Mutex
	out     io.Writer
	enabled [LEVEL_MAX]bool // enabled log level
	buf     []byte          // for accumulating text to write
}

// New creates a new Logger.
// The level variable sets the logger level. And the out variable sets the
// destination to which log data will be written.
func New(level Level, out io.Writer) *Logger {
	l := &Logger{out: out}
	for i := FATAL; i <= level && i < LEVEL_MAX; i++ {
		l.enabled[i] = true
	}
	return l
}

// The default logger
var std = New(DEBUG, os.Stderr)

// SetOutput sets the output destination for the standard logger.
func SetOutput(w io.Writer) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.out = w
}

// SetLevel sets the default logger level.
func SetLevel(level Level) {
	std.mu.Lock()
	defer std.mu.Unlock()
	for i := FATAL; i <= DEBUG; i++ {
		std.enabled[i] = i <= level
	}
}

// Output writes the output for a logging event.
// The level variable indicates the output message level. Note, only message
// level greater than or equals to logger level can be written. Calldepth is
// used to recover the PC and is provided for generality, although at the
// moment on all pre-defined paths it will be 2. The string s contains the text
// to print. A newline is appended if the last character of s is not already a
// newline.
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
	itoa(buf, t.Nanosecond()/1e6, 3)
	*buf = append(*buf, ": g="...)

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

// Fatalf is equivalent to Printf() for FATAL-level log.
func Fatalf(format string, v ...interface{}) {
	std.Output(FATAL, 2, fmt.Sprintf(format, v...))
}

// Fatal is equivalent to Print() for FATAL-level log.
func Fatal(v ...interface{}) {
	std.Output(FATAL, 2, fmt.Sprintln(v...))
}

// Warningf is equivalent to Printf() for WARNING-level log.
func Warningf(format string, v ...interface{}) {
	std.Output(WARNING, 2, fmt.Sprintf(format, v...))
}

// Waring is equivalent to Print() for WARING-level log.
func Warning(v ...interface{}) {
	std.Output(WARNING, 2, fmt.Sprintln(v...))
}

// Noticef is equivalent to Printf() for NOTICE-level log.
func Noticef(format string, v ...interface{}) {
	std.Output(NOTICE, 2, fmt.Sprintf(format, v...))
}

// Notice is equivalent to Print() for NOTICE-level log.
func Notice(v ...interface{}) {
	std.Output(NOTICE, 2, fmt.Sprintln(v...))
}

// Tracef is equivalent to Printf() for TRACE-level log.
func Tracef(format string, v ...interface{}) {
	std.Output(TRACE, 2, fmt.Sprintf(format, v...))
}

// Trace is equivalent to Print() for TRACE-level log.
func Trace(v ...interface{}) {
	std.Output(TRACE, 2, fmt.Sprintln(v...))
}

// Debugf is equivalent to Printf() for DEBUG-level log.
func Debugf(format string, v ...interface{}) {
	std.Output(DEBUG, 2, fmt.Sprintf(format, v...))
}

// Debug is equivalent to Print() for DEBUG-level log.
func Debug(v ...interface{}) {
	std.Output(DEBUG, 2, fmt.Sprintln(v...))
}
