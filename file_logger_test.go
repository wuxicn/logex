// -*- coding:utf-8; indent-tabs-mode:nil; -*-

// Copyright 2014, Wu Xi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logex

import (
    "fmt"
    "os"
	//"testing"
    //"time"
)

/* TODO: need fix this.
func TestClearLogs(t *testing.T) {
    KEEP_SECONDS = 1
    os.MkdirAll("./log/", 0755)
    writeFile("./log/test.log.1", "to be removed")
    writeFile("./log/test.log.2", "to be removed")
    d, _ := time.ParseDuration("1100ms")
    time.Sleep(d)
    writeFile("./log/test.log.3", "should be remained")

    err := SetUpFileLogger("./log", "test", func(now *time.Time, logfile *string) bool {
        return false
    })
    if err != nil {
        t.Fatal("set up error:", err)
    }

    d, _ = time.ParseDuration("1000ms")
    time.Sleep(d)

    // check if test.log.1 and test.log.2 are removed and test.log.3 is remained:
    if fileExists("./log/test.log.1") {
        t.Errorf("expect remove ./log/test.log.1, but it exists")
    }
    if fileExists("./log/test.log.2") {
        t.Errorf("expect remove ./log/test.log.2, but it exists")
    }
    if !fileExists("./log/test.log.3") {
        t.Errorf("expect remain ./log/test.log.3, but it not exists")
    }
}
*/

func writeFile(path, s string) {
    f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Printf("ERROR: open file %s for write failed: %v\n", path, err)
        os.Exit(255)
    }
    f.WriteString(s)
    f.Sync()
    f.Close()
    if !fileExists(path) {
        fmt.Printf("ERROR: file %s not exists\n", path)
        os.Exit(255)
    }
}

func fileExists(path string) bool {
    _, err := os.Stat(path)
    return err == nil || os.IsExist(err)
}

