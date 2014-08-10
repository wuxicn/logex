logex
=====

A leveled log pkg for go. Very easy to use and fully customizes.
Most importantly, it output goroutine-id in log line. output example:

    NOTICE: 08-10 15:15:03.365: 3: basic_usage.go:25: note the third field of log line is the goroutine-id


## Example

  * example.go:

    ``` go
    // The very basic usage of logex pkg
    package main

    import (
        "fmt"
        "os"
        "time"
    )

    import "github.com/wuxicn/logex"

    func main() {
        if err := SetUpFileLogger("./log", "test", nil); err != nil {
            fmt.Fprintln(os.Stderr, "init file logger failed\n")
            os.Exit(1)
        }
        logex.SetLevel(logex.NOTICE) // set log level to NOTICE,
                                     // so DEBUG and TRACE logs won't show
                                     // in log file

        logex.Debug("this message won't show")
        logex.Trace("this message won't show neither")
        logex.Notice("hi, pi is:", 3.1415926)
        logex.Warning("this is warning message")
        logex.Fatal("all logs will output to os.Stderr by default")
    }
    ```

## TODO



