// -*- coding:utf-8; indent-tabs-mode:nil; -*-

// Copyright 2014, Wu Xi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Output log to files.
// By default it writes to xxx.log file and saves xxx.log to xxx.log.<date> and
// opens an new xxx.log file when an new day begins. And it will also remove
// log files 7 days ago. Log file name, output dir and log file lifetime can be
// customized.
package logex

import (
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "path/filepath"
    "time"

    "github.com/mipearson/rfw"
)

var (
    LOGDIR_PERM  os.FileMode = 0755
    LOGFILE_PERM os.FileMode = 0644
    KEEP_SECONDS int = 86400 * 7 // keeps old log for 86400*7 seconds(7 days)
)

// LogChecker function returns true if need to close current log file and open
// an new one. now variable contains current time, and logfile contains
// absolute path of current log file.
type LogChecker func(now *time.Time, logfile *string) bool

type FileLogger struct {
    logdir  string
    prefix  string
    logfile string
    checker LogChecker
    writer  io.Writer
}

var logger FileLogger

func SetUpFileLogger(logdir, prefix string, checker LogChecker) error {
    absdir, err := filepath.Abs(logdir)
    if err != nil {
        return err
    }
    logger.prefix = prefix
    logger.logdir = absdir
    logger.logfile = fmt.Sprintf("%s/%s.log", absdir, prefix)
    if checker != nil {
        logger.checker = checker
    } else {
        logger.checker = DailyChecker
    }

    if err := os.MkdirAll(absdir, LOGDIR_PERM); err != nil {
        return err
    }
    logger.writer, err = rfw.Open(logger.logfile, LOGFILE_PERM)
    if err != nil {
        return err
    }

    SetOutput(logger.writer)

    go func() {
        layout := "2006-01-02_15-04-05"
        for logger.checker != nil {
            now := time.Now()
            if logger.checker(&now, &logger.logfile) {
                // move old log file
                savefile := logger.logfile + "." + now.Format(layout)
                if err := os.Rename(logger.logfile, savefile); err != nil {
                    fmt.Fprintln(os.Stderr, err)
                }
            }
            clearLogs()
            time.Sleep(time.Minute)
        }
    }()

    return nil
}

var daily = struct{
    layout string
    today  string
}{
    layout: "20060102",
    today: time.Now().Format("20060102"),
}

// DailyChecker returns true if now is the succeed day of daily.today
func DailyChecker(now *time.Time, logfile *string) bool {
    s := now.Format(daily.layout)
    if s != daily.today {
        daily.today = s
        return true
    }
    return false
}

// clearLogs clears log files older than (now - KEEP_SECONDS)
func clearLogs() {
    files, _ := ioutil.ReadDir(logger.logdir)
    l := len(logger.prefix)
    for _, f := range files {
        if len(f.Name()) > l && f.Name()[:l] == logger.prefix &&
            int(time.Since(f.ModTime()).Seconds()) > KEEP_SECONDS {
            Trace("remove old log:", f.Name())
            path := logger.logdir + "/" + f.Name()
            os.Remove(path)
        }
    }
}
