package logger

import (
	"io"
	"os"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var once sync.Once
var log zerolog.Logger

func Init(logLevel string) {
	parsedLogLevel, err := strconv.Atoi(logLevel)
	if err != nil {
		// default to INFO
		parsedLogLevel = int(zerolog.InfoLevel)
	}
	zerolog.SetGlobalLevel(zerolog.Level(parsedLogLevel))
}

func Get() zerolog.Logger {
	once.Do(func() {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.TimeFieldFormat = time.RFC3339Nano

		var output io.Writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}

		var gitRevision string

		buildInfo, ok := debug.ReadBuildInfo()
		if ok {
			for _, v := range buildInfo.Settings {
				if v.Key == "vcs.revision" {
					gitRevision = v.Value
					break
				}
			}
		}

		logLevel := zerolog.GlobalLevel()
		log = zerolog.New(output).
			Level(zerolog.Level(logLevel)).
			With().
			Timestamp().
			Str("git_revision", gitRevision).
			Str("go_version", buildInfo.GoVersion).
			Logger()
	})

	return log
}
