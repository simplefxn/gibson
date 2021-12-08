package logger

import (
	"encoding/json"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var zLogger zap.Logger

// RegisterLog ...
func RegisterLog() error {
	var err error
	zLogger, err = initLog()
	if err != nil {
		return errors.Wrap(err, "RegisterLog")
	}
	defer zLogger.Sync()
	zSugarlog := zLogger.Sugar()
	zSugarlog.Info()

	SetLogger(zSugarlog)
	return nil
}

// GetLogger get the configured logger
func GetLogger() *zap.Logger {
	return &zLogger
}

// initLog create logger
func initLog() (zap.Logger, error) {
	rawJSON := []byte(`{
	 "level": "info",
     "Development": true,
     "DisableCaller": false,
	 "encoding": "console",
	 "outputPaths": ["stdout", "./demo.log"],
	 "errorOutputPaths": ["stderr"],
	 "encoderConfig": {
		"timeKey":        "ts",
		"levelKey":       "level",
		"messageKey":     "msg",
        "nameKey":        "name",
		"stacktraceKey":  "stacktrace",
        "callerKey":      "caller",
		"lineEnding":     "\n\t",
        "timeEncoder":     "iso8601",
		"levelEncoder":    "lowercaseLevel",
        "durationEncoder": "stringDuration",
        "callerEncoder":   "shortCaller"
	 }
	}`)

	var cfg zap.Config
	var zLogger *zap.Logger
	var err error
	//standard configuration
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		return *zLogger, errors.Wrap(err, "Unmarshal")
	}

	zLogger, err = cfg.Build()
	if err != nil {
		return *zLogger, errors.Wrap(err, "cfg.Build()")
	}

	zLogger.Debug("logger construction succeeded")
	return *zLogger, nil
}
