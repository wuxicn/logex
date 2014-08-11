// -*- coding:utf-8; indent-tabs-mode:nil; -*-

// Copyright 2014, Wu Xi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logex

import (
    "errors"
    "io/ioutil"
    "fmt"
    "os"
    "runtime"
    "strings"
    "strconv"
	"testing"
    "time"
)

func TestNormalFileLogger(t *testing.T) {
    f := "test_normal_file_logger"
    os.Remove("./log/" + f + ".log")

    err := SetUpFileLogger("./log", f, func(now *time.Time, logfile *string) bool {
        return false
    })
    if err != nil {
        t.Fatal("set up error:", err)
    }

    _, _, line, _ := runtime.Caller(0)
    Notice("abc")
    if err := checkFile(f, NOTICE, line + 1, "abc\n"); err != nil {
        t.Fatal(err)
    }
}

func checkFile(f string, level Level, lineno int, msg string) error {
    b, err := ioutil.ReadFile("./log/" + f + ".log")
    if err != nil {
        return err
    }

    a := strings.SplitN(string(b), ": ", 5)
    if len(a) != 5 {
        return errors.New("wrong log line format")
    }

    if levelStr[level] != a[0] {
        return errors.New(fmt.Sprintf("expect level=%q but actually is %q",
            levelStr[level], a[0]))
    }

    gid := strconv.Itoa(int(goid()))
    if gid != a[2] {
        return errors.New(fmt.Sprintf("expect gid=%q but actually is %q",
            gid, a[2]))
    }

    s := fmt.Sprintf("file_logger_test.go:%d", lineno)
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



