package main

import "github.com/gaukas/passthru/internal/logger"

func main() {
	logger.InitLogger("test.log", true, logger.LOG_DEBUG)

	logger.Debugf("Debugf message")
	logger.Infof("Infof message")
	logger.Warnf("Warnf message")
	logger.Errorf("Errorf message")
	logger.Fatalf("Fatalf message")

	logger.Errorf("Failed.")
}
