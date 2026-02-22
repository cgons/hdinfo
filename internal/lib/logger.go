package lib

import (
	"sync"

	"go.uber.org/zap"
)

var logger *zap.SugaredLogger
var once sync.Once

// InitLogger initializes a thread-safe singleton logger
func InitLogger() {
	once.Do(func() {
		var err error
		zapLogger, err := zap.NewDevelopment()
		logger = zapLogger.Sugar()
		defer zapLogger.Sync()
		if err != nil {
			panic(err)
		}
	})
}

func GetLogger() *zap.SugaredLogger {
	return logger
}
