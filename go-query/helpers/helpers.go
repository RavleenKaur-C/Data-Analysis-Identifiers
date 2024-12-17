package helpers

import (
	"os"
	"regexp"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitializeLogger() *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	level := zapcore.DebugLevel
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		level,
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	return logger
}
func IsSHAOrUUID(s string) bool {

	uuidRegex := regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`)
	shaRegex := regexp.MustCompile(`^[a-fA-F0-9]+$`)

	validSHALengths := map[int]string{
		40:  "SHA-1",
		48:  "SHA-224",
		56:  "SHA-512/224, SHA3-224",
		64:  "SHA-256, SHA3-256, SHAKE256",
		96:  "SHA-384, SHA3-384",
		128: "SHA-512, SHA3-512",
		32:  "SHAKE128",
	}

	if len(s) == 36 && uuidRegex.MatchString(s) {
		return true
	}

	if shaRegex.MatchString(s) {
		if _, exists := validSHALengths[len(s)]; exists {
			return true
		}
	}

	return false
}
