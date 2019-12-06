package golog

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/widaT/go-log/log4go"
)

//Logger 单例
var Logger log4go.Logger
var mutex sync.Mutex
var levelMapping = map[string]log4go.LevelType{
	"DEBUG":    log4go.DEBUG,
	"TRACE":    log4go.TRACE,
	"INFO":     log4go.INFO,
	"WARNING":  log4go.WARNING,
	"ERROR":    log4go.ERROR,
	"CRITICAL": log4go.CRITICAL,
}

func createDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func createFile(prefix, logDir string, isErrLog bool) string {
	strings.TrimSuffix(logDir, "/")
	var fileName string
	if isErrLog {
		fileName = filepath.Join(logDir, prefix+".wf.log")
	} else {
		fileName = filepath.Join(logDir, prefix+".log")
	}
	return fileName
}

// Init 初始化
//
// PARAMS:
//   - prefix: 日志前缀
//   - levelStr: "DEBUG", "TRACE", "INFO", "WARNING", "ERROR", "CRITICAL"
//   - logDir: 日志目录，没有的话会自动创建
//   - filterWf: 是否过滤WARNING以上的日志到另外一个文件
//   - when:
//       "M", minute
//       "H", hour
//       "D", day
//       "MIDNIGHT", roll over at midnight
//       "NEXTHOUR", roll over at sharp next hour
//   - backupCount: If backupCount is > 0, when rollover is done, no more than
//       backupCount files are kept - the oldest ones are deleted.
func Init(prefix string, levelStr string, logDir string, filterWf bool,
	when string, backupCount int) error {
	mutex.Lock()
	defer mutex.Unlock()

	if Logger != nil {
		return errors.New("Initialized Already")
	}

	var err error
	Logger, err = GetLogger(prefix, levelStr, logDir, filterWf, when, backupCount)
	if err != nil {
		return err
	}
	return nil
}

// GetLogger creates logger
//
// PARAMS:
//   - prefix: 日志前缀
//   - levelStr: "DEBUG", "TRACE", "INFO", "WARNING", "ERROR", "CRITICAL"
//   - logDir: 日志目录，没有的话会自动创建
//   - filterWf: 是否过滤WARNING以上的日志到另外一个文件
//   - when:
//       "M", minute
//       "H", hour
//       "D", day
//       "MIDNIGHT", roll over at midnight
//       "NEXTHOUR", roll over at sharp next hour
//   - backupCount: If backupCount is > 0, when rollover is done, no more than
//       backupCount files are kept - the oldest ones are deleted.
func GetLogger(prefix string, levelStr string, logDir string, filterWf bool,
	when string, backupCount int) (log4go.Logger, error) {
	if !log4go.WhenIsValid(when) {
		return nil, fmt.Errorf("invalid value of when: %s", when)
	}

	if err := createDir(logDir); err != nil {
		return nil, fmt.Errorf("createDir error %s", err.Error())
	}

	levelStr = strings.ToUpper(levelStr)
	level, found := levelMapping[levelStr]
	if !found {
		return nil, fmt.Errorf("invalid levelStr %s", levelStr)
	}

	logger := make(log4go.Logger)
	/* 	if withStdOut {
		logger.AddFilter("stdout", level, log4go.NewConsoleLogWriter())
	} */

	fileName := createFile(prefix, logDir, false)
	logWriter := log4go.NewTimeFileLogWriter(fileName, when, backupCount)
	if logWriter == nil {
		return nil, fmt.Errorf("log4go.NewTimeFileLogWriter error %s", fileName)
	}
	logWriter.SetFormat(log4go.LogFormat)
	logger.AddFilter("log", level, logWriter)

	if filterWf {
		fileNameWf := createFile(prefix, logDir, true)
		logWriter = log4go.NewTimeFileLogWriter(fileNameWf, when, backupCount)
		if logWriter == nil {
			return nil, fmt.Errorf("log4go.NewTimeFileLogWriter error %s", fileNameWf)
		}
		logWriter.SetFormat(log4go.LogFormat)
		logger.AddFilter("log_wf", log4go.WARNING, logWriter)
	}

	return logger, nil
}
