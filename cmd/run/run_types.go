package run

import (
	"io"

	"samhofi.us/x/keybase/v2"
	"samhofi.us/x/keybase/v2/types/chat1"
)

type bot struct {
	k        *keybase.Keybase
	handlers keybase.Handlers
	opts     keybase.RunOptions
	config   botConfig
}

type botConfig struct {
	debug        bool             // Whether to enable debugging
	logConv      *chat1.ConvIDStr // ConversationID to send log messages to. Set to nil to disable
	settingsTeam *string          // Team to use for storing team settings. Set to nil to use implicit self-team (the bot user's conversation with itself)
	stdout       io.Writer
	stderr       io.Writer
}

type logLevel int

const (
	logLevelUnknown logLevel = iota
	logLevelInfo
	logLevelDebug
	logLevelError
)

var logLevelStringMap = map[string]logLevel{
	"UNKNOWN": logLevelUnknown,
	"INFO":    logLevelInfo,
	"DEBUG":   logLevelDebug,
	"ERROR":   logLevelError,
}

var logLevelStringRevMap = map[logLevel]string{
	logLevelUnknown: "UNKNOWN",
	logLevelInfo:    "INFO",
	logLevelDebug:   "DEBUG",
	logLevelError:   "ERROR",
}

func (l logLevel) String() string {
	return logLevelStringRevMap[l]
}
