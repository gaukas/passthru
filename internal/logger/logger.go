package logger

import (
	"log"
	"os"
	"sync"
)

const (
	LOG_DEBUG uint8 = iota
	LOG_INFO
	LOG_WARN
	LOG_ERROR
	LOG_FATAL
	LOG_NOLOG
)

var (
	fileLogger *log.Logger

	loggerInitOnce sync.Once

	logFile *os.File

	_Debugf     func(format string, v ...interface{}) = func(string, ...interface{}) {}
	_Infof      func(format string, v ...interface{}) = func(string, ...interface{}) {}
	_Warnf      func(format string, v ...interface{}) = func(string, ...interface{}) {}
	_Errorf     func(format string, v ...interface{}) = func(string, ...interface{}) {}
	_Fatalf     func(format string, v ...interface{}) = func(string, ...interface{}) {}
	_FileDebugf func(format string, v ...interface{}) = func(string, ...interface{}) {}
	_FileInfof  func(format string, v ...interface{}) = func(string, ...interface{}) {}
	_FileWarnf  func(format string, v ...interface{}) = func(string, ...interface{}) {}
	_FileErrorf func(format string, v ...interface{}) = func(string, ...interface{}) {}
	_FileFatalf func(format string, v ...interface{}) = func(string, ...interface{}) {}
)

func InitLogger(filename string, stderr bool, level uint8) {
	loggerInitOnce.Do(func() {
		if filename != "" {
			var err error
			logFile, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatalf("can't initialize log file, os.OpenFile: %v", err)
			}
			fileLogger = log.New(logFile, "", log.LstdFlags)
		}

		setLoggingLevel(level, stderr)
	})
}

// Inmutable binding to log functions
func Debugf(format string, v ...interface{}) {
	_Debugf(format, v...)
	_FileDebugf(format, v...)
}

func Infof(format string, v ...interface{}) {
	_Infof(format, v...)
	_FileInfof(format, v...)
}

func Warnf(format string, v ...interface{}) {
	_Warnf(format, v...)
	_FileWarnf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	_Errorf(format, v...)
	_FileErrorf(format, v...)
}

func Fatalf(format string, v ...interface{}) {
	_Fatalf(format, v...)
	_FileFatalf(format, v...)
	os.Exit(1)
}

func setLoggingLevel(level uint8, stderr bool) {
	if stderr {
		if level <= LOG_FATAL {
			_Fatalf = func(format string, v ...interface{}) {
				log.Printf("["+red+"FATAL"+reset+"] "+red+format+reset, v...)
			}
		}
		if level <= LOG_ERROR {
			_Errorf = func(format string, v ...interface{}) {
				log.Printf("["+red+"ERROR"+reset+"] "+format, v...)
			}
		}
		if level <= LOG_WARN {
			_Warnf = func(format string, v ...interface{}) {
				log.Printf("["+yellow+"WARN"+reset+"] "+format, v...)
			}
		}
		if level <= LOG_INFO {
			_Infof = func(format string, v ...interface{}) {
				log.Printf("["+green+"INFO"+reset+"] "+format, v...)
			}
		}
		if level <= LOG_DEBUG {
			_Debugf = func(format string, v ...interface{}) {
				log.Printf("["+blue+"DEBUG"+reset+"] "+format, v...)
			}
		}
	}

	if fileLogger != nil {
		if level <= LOG_FATAL {
			_FileFatalf = func(format string, v ...interface{}) {
				fileLogger.Printf("[FATAL] "+format, v...)
			}
		}
		if level <= LOG_ERROR {
			_FileErrorf = func(format string, v ...interface{}) {
				fileLogger.Printf("[ERROR] "+format, v...)
			}
		}
		if level <= LOG_WARN {
			_FileWarnf = func(format string, v ...interface{}) {
				fileLogger.Printf("[WARN] "+format, v...)
			}
		}
		if level <= LOG_INFO {
			_FileInfof = func(format string, v ...interface{}) {
				fileLogger.Printf("[INFO] "+format, v...)
			}
		}
		if level <= LOG_DEBUG {
			_FileDebugf = func(format string, v ...interface{}) {
				fileLogger.Printf("[DEBUG] "+format, v...)
			}
		}
	}

}
