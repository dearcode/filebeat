package lb

import (
	"testing"

	"github.com/dearcode/libbeat/common"
	"github.com/dearcode/libbeat/logp"
	"github.com/dearcode/libbeat/outputs"
)

var (
	testNoOpts     = outputs.Options{}
	testGuaranteed = outputs.Options{Guaranteed: true}

	testEvent = common.MapStr{
		"msg": "hello world",
	}
)

func enableLogging(selectors []string) {
	if testing.Verbose() {
		logp.LogInit(logp.LOG_DEBUG, "", false, true, selectors)
	}
}
