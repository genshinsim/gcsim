package monster

import (
	"os"
	"testing"

	"github.com/genshinsim/gsim/pkg/core"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	config := zap.NewDevelopmentConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	config.EncoderConfig.TimeKey = ""
	log, _ := config.Build(zap.AddCallerSkip(1))
	logger = log.Sugar()
	os.Exit(m.Run())
}

func durApproxEqual(expect, result, tol core.Durability) bool {
	if expect > result {
		return expect-result < tol
	}
	return result-expect < tol
}

func floatApproxEqual(expect, result, tol float64) bool {
	if expect > result {
		return expect-result < tol
	}
	return result-expect < tol
}

func expect(msg string, a interface{}, b interface{}) {
	logger.Infow(msg, "expected", a, "got", b)
}

func defaultResMap() map[core.EleType]float64 {
	res := make(map[core.EleType]float64)

	res[core.Electro] = 0.1
	res[core.Pyro] = 0.1
	res[core.Anemo] = 0.1
	res[core.Cryo] = 0.1
	res[core.Frozen] = 0.1
	res[core.Hydro] = 0.1
	res[core.Dendro] = 0.1
	res[core.Geo] = 0.1
	res[core.Physical] = 0.1

	return res
}
