package logger

import (
	"io"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var L *zap.Logger

// LogEntry represents a single log entry returned by GetRecentLogs.
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

// ringBuffer is a thread-safe fixed-size ring buffer for LogEntry.
type ringBuffer struct {
	mu     sync.Mutex
	buf    []LogEntry
	size   int
	cursor int
	full   bool
}

func newRingBuffer(size int) *ringBuffer {
	return &ringBuffer{
		buf:  make([]LogEntry, size),
		size: size,
	}
}

func (rb *ringBuffer) Push(entry LogEntry) {
	rb.mu.Lock()
	rb.buf[rb.cursor] = entry
	rb.cursor++
	if rb.cursor >= rb.size {
		rb.cursor = 0
		rb.full = true
	}
	rb.mu.Unlock()
}

// Recent returns the n most recent entries, in chronological order.
func (rb *ringBuffer) Recent(n int) []LogEntry {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if n <= 0 {
		return []LogEntry{}
	}
	if n > rb.size {
		n = rb.size
	}

	var total int
	if rb.full {
		total = rb.size
	} else {
		total = rb.cursor
	}
	if n > total {
		n = total
	}
	if n == 0 {
		return []LogEntry{}
	}

	result := make([]LogEntry, n)
	if rb.full {
		start := (rb.cursor - n + rb.size) % rb.size
		for i := 0; i < n; i++ {
			result[i] = rb.buf[(start+i)%rb.size]
		}
	} else {
		copy(result, rb.buf[rb.cursor-n:rb.cursor])
	}
	return result
}

var logBuf = newRingBuffer(1000)

// GetRecentLogs returns the n most recent log entries, in chronological order.
func GetRecentLogs(n int) []LogEntry {
	return logBuf.Recent(n)
}

// ringBufferCore is a zapcore.Core that captures log entries into the ring buffer
// while delegating actual output to the wrapped core.
type ringBufferCore struct {
	zapcore.Core
	rb *ringBuffer
}

func (c *ringBufferCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	c.rb.Push(LogEntry{
		Timestamp: entry.Time.Format(time.RFC3339Nano),
		Level:     entry.Level.String(),
		Message:   entry.Message,
	})
	return c.Core.Write(entry, fields)
}

func (c *ringBufferCore) Check(entry zapcore.Entry, checkedEntry *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return checkedEntry.AddCore(entry, c)
	}
	return checkedEntry
}

func (c *ringBufferCore) With(fields []zapcore.Field) zapcore.Core {
	return &ringBufferCore{
		Core: c.Core.With(fields),
		rb:   c.rb,
	}
}

func newRingBufferCore(primary zapcore.Core) zapcore.Core {
	return &ringBufferCore{
		Core: primary,
		rb:   logBuf,
	}
}

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
	// Wrap with ring buffer capture.
	L = L.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return newRingBufferCore(core)
	}))
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
	primary := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(os.Stdout),
		zap.NewAtomicLevelAt(zap.DebugLevel),
	)
	L = zap.New(newRingBufferCore(primary), zap.AddCallerSkip(1))
	if env != "production" {
		L = L.WithOptions(zap.Development())
	}
}
