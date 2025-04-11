package logger

import (
    "fmt"
    "io"
    "log"
    "os"
    "time"
)

// InitLog - создаёт файл для лога и инициализирует глобальный логгер
func InitLog(logPath string, enableConsole bool) {

    if logPath == "" {
        logPath == "logs" + string(os.PathSeparator)
    }
    err := os.MkdirAll(logPath, 0755)
    if err != nil {
        log.Fatalf("Error creating directories for logs: %v", err)
    }

    logFileName:= logPath + fmt.Sprintf("log_%s.log", time.Now().Format("20060101_150405"))
    logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("Ошибка открытия лог-файла: %v", err)
    }
    var logger io.Writer
    if enableConsole {
        logger = io.MultiWriter(os.Stdout, logFile)
    } else {
        logger = logFile
    }
    log.SetOutput(logger)
    log.SetFlags(log.Ldate | log.Lmicroseconds)
    log.Println("Logger was initializing. Log file:", logFileName)
}
