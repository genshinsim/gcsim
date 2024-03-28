package optimization

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newLogger(verbose bool) *zap.SugaredLogger {
	// Start logger
	zapcfg := zap.NewDevelopmentConfig()
	zapcfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	zapcfg.EncoderConfig.CallerKey = ""
	zapcfg.EncoderConfig.StacktraceKey = ""
	zapcfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	if verbose {
		zapcfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}

	logger, _ := zapcfg.Build()
	defer logger.Sync()
	sugarLog := logger.Sugar()

	return sugarLog
}
