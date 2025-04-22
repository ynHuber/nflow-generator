package main

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

type ZeroLogger struct {
}

func (l *ZeroLogger) Debug(args ...interface{}) {
	log.Debug().Msg(fmt.Sprint(args...))
}

func (l *ZeroLogger) Debugf(format string, args ...interface{}) {
	log.Debug().Msgf(format, args...)
}

func (l *ZeroLogger) Debugln(args ...interface{}) {
	log.Debug().Msg(fmt.Sprintln(args...))
}

func (l *ZeroLogger) Error(args ...interface{}) {
	log.Error().Msg(fmt.Sprint(args...))
}

func (l *ZeroLogger) Errorf(format string, args ...interface{}) {
	log.Error().Msgf(format, args...)
}

func (l *ZeroLogger) Errorln(args ...interface{}) {
	log.Error().Msg(fmt.Sprintln(args...))
}

func (l *ZeroLogger) Fatal(args ...interface{}) {
	log.Fatal().Msg(fmt.Sprint(args...))
}

func (l *ZeroLogger) Fatalf(format string, args ...interface{}) {
	log.Fatal().Msgf(format, args...)
}

func (l *ZeroLogger) Fatalln(args ...interface{}) {
	log.Fatal().Msg(fmt.Sprintln(args...))
}

func (l *ZeroLogger) Info(args ...interface{}) {
	log.Info().Msg(fmt.Sprint(args...))
}

func (l *ZeroLogger) Infof(format string, args ...interface{}) {
	log.Info().Msgf(format, args...)
}

func (l *ZeroLogger) Infoln(args ...interface{}) {
	log.Info().Msg(fmt.Sprintln(args...))
}

func (l *ZeroLogger) Trace(args ...interface{}) {
	log.Trace().Msg(fmt.Sprint(args...))
}

func (l *ZeroLogger) Tracef(format string, args ...interface{}) {
	log.Trace().Msgf(format, args...)
}

func (l *ZeroLogger) Traceln(args ...interface{}) {
	log.Trace().Msg(fmt.Sprintln(args...))
}

func (l *ZeroLogger) Warn(args ...interface{}) {
	log.Warn().Msg(fmt.Sprint(args...))
}

func (l *ZeroLogger) Warnf(format string, args ...interface{}) {
	log.Warn().Msgf(format, args...)
}

func (l *ZeroLogger) Warning(args ...interface{}) {
	log.Warn().Msg(fmt.Sprint(args...))
}

func (l *ZeroLogger) Warningf(format string, args ...interface{}) {
	log.Warn().Msgf(format, args...)
}

func (l *ZeroLogger) Warningln(args ...interface{}) {
	log.Warn().Msg(fmt.Sprintln(args...))
}

func (l *ZeroLogger) Warnln(args ...interface{}) {
	log.Warn().Msg(fmt.Sprintln(args...))
}
