package jujuclient

import (
	"context"
	"fmt"

	"github.com/juju/juju/core/logger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type LoggerWrapper struct {
	logger zerolog.Logger
}

func NewLoggerWrapper() logger.Logger {
	return &LoggerWrapper{
		logger: log.Logger,
	}
}

func (l *LoggerWrapper) getMessage(msg string, args ...any) string {
	return fmt.Sprintf(msg, args...)
}

func (l *LoggerWrapper) Criticalf(ctx context.Context, msg string, args ...any) {
	log.Ctx(ctx).Fatal().Msg(l.getMessage(msg, args))
}

func (l *LoggerWrapper) Errorf(ctx context.Context, msg string, args ...any) {
	log.Ctx(ctx).Error().Msg(l.getMessage(msg, args))
}

func (l *LoggerWrapper) Warningf(ctx context.Context, msg string, args ...any) {
	log.Ctx(ctx).Warn().Msg(l.getMessage(msg, args))
}

func (l *LoggerWrapper) Infof(ctx context.Context, msg string, args ...any) {
	log.Ctx(ctx).Info().Msg(l.getMessage(msg, args))
}

func (l *LoggerWrapper) Debugf(ctx context.Context, msg string, args ...any) {
	log.Ctx(ctx).Debug().Msg(l.getMessage(msg, args))
}

func (l *LoggerWrapper) Tracef(ctx context.Context, msg string, args ...any) {
	log.Ctx(ctx).Trace().Msg(l.getMessage(msg, args))
}

func (l *LoggerWrapper) Logf(ctx context.Context, level logger.Level, labels logger.Labels, format string, args ...any) {
	log.Ctx(ctx).WithLevel(l.logger.GetLevel()).Msgf(format, args)
}

func (l *LoggerWrapper) IsLevelEnabled(level logger.Level) bool {
	if int(level) < int(l.logger.GetLevel()) {
		return false
	}
	return true
}

func (l *LoggerWrapper) Child(name string, tags ...string) logger.Logger {
	return &LoggerWrapper{
		logger: log.Level(log.Logger.GetLevel()),
	}
}

func (l *LoggerWrapper) GetChildByName(name string) logger.Logger {
	return l
}
