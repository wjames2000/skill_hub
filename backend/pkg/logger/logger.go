package logger

import (
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var L *zap.Logger

func Init(env string) {
	var cfg zap.Config
	if env == "production" {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}
	cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	var err error
	L, err = cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic("logger init: " + err.Error())
	}
}

func NewGinWriter() io.Writer {
	r, w := io.Pipe()
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := r.Read(buf)
			if err != nil {
				return
			}
			L.Info(string(buf[:n]))
		}
	}()
	return w
}

func Info(msg string, fields ...zap.Field)  { L.Info(msg, fields...) }
func Warn(msg string, fields ...zap.Field)  { L.Warn(msg, fields...) }
func Error(msg string, fields ...zap.Field) { L.Error(msg, fields...) }
func Debug(msg string, fields ...zap.Field) { L.Debug(msg, fields...) }
func Fatal(msg string, fields ...zap.Field) { L.Fatal(msg, fields...) }
func Sync()                                 { _ = L.Sync() }

func String(key, val string) zap.Field                 { return zap.String(key, val) }
func Int(key string, val int) zap.Field                { return zap.Int(key, val) }
func Int64(key string, val int64) zap.Field            { return zap.Int64(key, val) }
func Bool(key string, val bool) zap.Field              { return zap.Bool(key, val) }
func Any(key string, val interface{}) zap.Field        { return zap.Any(key, val) }
func Duration(key string, val time.Duration) zap.Field { return zap.Duration(key, val) }
func ErrorField(err error) zap.Field                   { return zap.Error(err) }

func init() {
	env := os.Getenv("SKILL_HUB_ENV")
	if env == "" {
		env = "development"
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(os.Stdout),
		zap.NewAtomicLevelAt(zap.DebugLevel),
	)
	L = zap.New(core, zap.AddCallerSkip(1))
	if env != "production" {
		L = L.WithOptions(zap.Development())
	}
}
