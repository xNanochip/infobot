package run

import (
	"fmt"
	"strings"
	"time"

	"samhofi.us/x/infobot/pkg/utils"
)

func (b *bot) write_log(level logLevel, s string, a ...interface{}) {
	var (
		now        = time.Now().UTC()
		timeFormat = "02Jan2006 15:04:05"
	)

	if b.config.json {
		item := logItem{
			Time:       now.Unix(),
			LogLevel:   level.String(),
			LogMessage: fmt.Sprintf(s, a...),
		}
		fmt.Fprintln(b.config.stdout, utils.ToJson(item))
		return
	}

	a = append([]interface{}{strings.ToUpper(now.Format(timeFormat)), level}, a...)
	fmt.Fprintf(b.config.stdout, "[%v] %v: "+s+"\n", a...)
}

func (b *bot) log_error(s string, a ...interface{}) {
	b.write_log(logLevelError, s, a...)
}

func (b *bot) log_info(s string, a ...interface{}) {
	b.write_log(logLevelInfo, s, a...)
}

func (b *bot) log_debug(s string, a ...interface{}) {
	if !b.config.debug {
		return
	}
	b.write_log(logLevelDebug, s, a...)
}
