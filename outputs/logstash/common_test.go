package logstash

import (
	"testing"

	"github.com/dearcode/libbeat/logp"
)

func enableLogging(selectors []string) {
	if testing.Verbose() {
		logp.LogInit(logp.LOG_DEBUG, "", false, true, selectors)
	}
}
