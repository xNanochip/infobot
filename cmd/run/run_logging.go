package run

import (
	"fmt"
	"strings"
	"time"
)

func (b *bot) write_log(level logLevel, s string, a ...interface{}) {
	var currentTime = strings.ToUpper(time.Now().UTC().Format("02Jan2006 15:04:05"))
	a = append([]interface{}{currentTime, level}, a...)
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
