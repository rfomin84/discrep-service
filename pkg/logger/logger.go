package logger

import (
	zerolog "github.com/rs/zerolog"
	"os"
)

var log *zerolog.Logger

func init() {
	if log == nil {
		zero := zerolog.New(os.Stdout)
		log = &zero
	}
}

func SetLevel(level zerolog.Level) {
	log.Level(level)
}

func Info(message string) {
	log.Info().Timestamp().Msg(message)
}

func Debug(message string) {
	log.Debug().Timestamp().Msg(message)
}

func Warning(message string) {
	log.Warn().Timestamp().Msg(message)
}

func Error(message string) {
	log.Error().Timestamp().Msg(message)
}
