### Introduction ###

This is a log manager tool,it base on https://github.com/sirupsen/logrus

### How to use ###

** basic usage **

```
package main
import "bitbucket.org/bindolabs/golib/logger"

func main(){
    logger.Info("somethins")
    logger.Infof("input sth %s", "sth")
}
```


** change logger target to some file

```
    filehandle ,_ := logger.CreateFileLogger("logpath.log", logger.FileLoggerNone)
    LogManager["default"] = &filehandle.Logger
```

** logger.FileLoggerNone ** means donot use logger built-in log rotation tool


if you want rotation log file with Hour or Day use logger.FileLoggerHour/logger.FileLoggerDay

