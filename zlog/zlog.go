package zlog

import "go.uber.org/zap"

func init() {
	zlog, err := zap.NewDevelopment()
	if err != nil {
		panic(" init the zlog fail")
	}
	zap.ReplaceGlobals(zlog)
}

func Panic(args ...any) {
	zap.S().Panic(args)
}

func Error(args ...any) {
	zap.S().Error(args)
}

func Warn(args ...any) {
	zap.S().Warn(args)
}

func Info(args ...any) {
	zap.S().Info(args)
}

func Debug(args ...any) {
	zap.S().Debug(args)
}
