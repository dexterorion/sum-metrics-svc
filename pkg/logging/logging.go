package logging

import (
	"os"
	"time"

	"github.com/blendle/zapdriver"
	env_vars "github.com/dexterorion/mao-backend/helpers/envvars"
	"github.com/onsi/gomega/gstruct/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	root        *zap.Logger
	factory     *zap.SugaredLogger
	main        *zap.SugaredLogger
	toflush     = []*zap.SugaredLogger{}
	devMode     = os.Getenv(env_vars.IS_DEV_KEY) == env_vars.IS_DEV_VALUE
	enableDebug = os.Getenv(env_vars.DEBUG_MODE_KEY) == env_vars.IS_DEBUG_VALUE
)

func AddLabel(name string, value string) {
	if root != nil {
		root.With(zapdriver.Label(name, value))
	}
}

func Init(service string) *zap.SugaredLogger {
	var config zap.Config

	if devMode {
		config = buildDevelopmentConfig()
	} else {
		config = buildProductionConfig()
	}

	// hooks := zap.Hooks(PromHook)

	if factory == nil {
		var err error
		if root, err = config.Build(getStackLevel(), zapdriver.WrapCore() /*, hooks*/); err == nil {
			if service != "" {
				root.
					With(zapdriver.Label("service", service)).
					With(zapdriver.Label("k8s-pod/app", service))
			}

			factory = root.Sugar()

			main = New("main")

			return main
		} else {
			panic(errors.Nest("could not initialize zap logger", err))
		}
	} else {
		root.
			With(zapdriver.Label("service", service)).
			With(zapdriver.Label("k8s-pod/app", service))

		return main
	}
}

func buildDevelopmentConfig() zap.Config {
	encoderConfig := zapdriver.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeTime = func(t time.Time, e zapcore.PrimitiveArrayEncoder) {

		e.AppendString(t.Format("15:04:05.000"))
	}
	encoderConfig.EncodeName = func(loggerName string, enc zapcore.PrimitiveArrayEncoder) {
		if len(loggerName) >= 10 {
			enc.AppendString(loggerName[:9])
		} else {
			enc.AppendString(PadRight(loggerName, " ", 10))
		}
	}

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
		DisableCaller:    true,
	}

	return config
}

func buildProductionConfig() zap.Config {

	level := zap.NewAtomicLevelAt(zap.InfoLevel)
	disableCaller := true

	if enableDebug {
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
		disableCaller = false
	}

	config := zap.Config{
		Level:            level,
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    zapdriver.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
		DisableCaller:    disableCaller,
	}

	return config
}

func PadRight(str, pad string, lenght int) string {
	for {
		str += pad
		if len(str) > lenght {
			return str[0:lenght]
		}
	}
}

func New(name string) *zap.SugaredLogger {
	if factory == nil {
		Init("")
		//factory.Warn("Logging was not initialized. If this isn't a test run, maybe you should review the program startup.")
	}

	l := factory.Named(name)

	toflush = append(toflush, l)

	return l
}

func Flush() {
	main.Info("flushing logs...")

	for _, v := range toflush {
		v.Sync()
	}

	factory.Sync()
}

func getStackLevel() zap.Option {
	if devMode {
		return zap.AddStacktrace(zapcore.ErrorLevel)
	} else {
		return zap.AddStacktrace(zapcore.FatalLevel)
	}
}
