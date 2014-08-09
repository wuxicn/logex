// -*- coding:utf-8; indent-tabs-mode:nil; -*-

// Copyright 2014, Wu Xi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logex

import (
    "bytes"
    "errors"
	"fmt"
    "runtime"
    "strings"
    "strconv"
	"testing"
)

func TestNormalOutput(t *testing.T) {
	buf := new(bytes.Buffer)
	SetOutput(buf)
	SetLevel(DEBUG)

    _, _, line, _ := runtime.Caller(0)
	Fatal("hello", "world")
    if err := check(buf, FATAL, line + 1, "hello world\n"); err != nil {
        t.Error(err)
    }

    buf.Reset()
    _, _, line2, _ := runtime.Caller(0)
    Debugf("abc %d xyz\n", 123)
    if err := check(buf, DEBUG, line2 + 1, "abc 123 xyz\n"); err != nil {
        t.Error(err)
    }
}

func TestLogLevel(t *testing.T) {
	buf := new(bytes.Buffer)
	SetOutput(buf)
	SetLevel(NONE)

    // will not write any bytes in buf
	Fatal("hello", "world")
    if buf.Len() > 0 {
        t.Errorf("expect writes nothing, but actually writes %q", buf.String())
    }

    // will output nothing but fatal and warning:
    SetLevel(WARNING)

    buf.Reset()
    _, _, line, _ := runtime.Caller(0)
	Fatal("abc")
    if err := check(buf, FATAL, line + 1, "abc\n"); err != nil {
        t.Error(err)
    }

    buf.Reset()
    _, _, line2, _ := runtime.Caller(0)
	Warning("def")
    if err := check(buf, WARNING, line2 + 1, "def\n"); err != nil {
        t.Error(err)
    }

    buf.Reset()
	Notice("abc")
    if buf.Len() > 0 {
        t.Errorf("expect writes nothing, but actually writes %q", buf.String())
    }

    buf.Reset()
	Trace("abc")
    if buf.Len() > 0 {
        t.Errorf("expect writes nothing, but actually writes %q", buf.String())
    }

    buf.Reset()
	Debug("abc")
    if buf.Len() > 0 {
        t.Errorf("expect writes nothing, but actually writes %q", buf.String())
    }
}

// check log output line
// line eg: "FATAL", "08-09 17:03:11.994", "3", "log_test.go:24", "hello world\n"
func check(b *bytes.Buffer, level Level, lineno int, msg string) error {
    a := strings.SplitN(b.String(), ": ", 5)
    if len(a) != 5 {
        return errors.New("wrong log line format")
    }

    if levelStr[level] != a[0] {
        return errors.New(fmt.Sprintf("expect level=%q but actually is %q",
            levelStr[level], a[0]))
    }

    gid := strconv.Itoa(int(goid()))
    if gid != a[2] {
        return errors.New(fmt.Sprintf("expect goid=%q but actually is %q",
            gid, a[2]))
    }

    s := fmt.Sprintf("log_test.go:%d", lineno)
    if s != a[3] {
        return errors.New(fmt.Sprintf("expect file:line=%q but actually is %q",
            s, a[3]))
    }

    if msg != a[4] {
        return errors.New(fmt.Sprintf("expect msg=%q but actually is %q",
            msg, a[4]))
    }

    return nil
}

