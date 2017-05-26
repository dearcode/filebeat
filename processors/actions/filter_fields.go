package actions

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/dearcode/libbeat/common"
	"github.com/dearcode/libbeat/processors"
)

type TimeValue struct {
	s  int64
	ms time.Duration
	sync.Mutex
}

type filterFields struct {
	Topic      string
	Key        string
	Regexp     *regexp.Regexp
	Names      []string
	DateLayout string
	timeValue  *TimeValue
}

func (tv *TimeValue) timestamp(t time.Time) string {
	tv.Lock()
	defer tv.Unlock()

	if t.Unix() != tv.s {
		tv.ms = 0
		tv.s = t.Unix()
	}
	tv.ms++
	return t.Add(tv.ms * time.Millisecond).UTC().Format(common.TsLayout)
}

func init() {
	processors.RegisterPlugin("filter_fields", configChecked(newFilterFields, requireFields("topic", "regexp", "names"), allowedFields("topic", "key", "regexp", "when", "names", "date_layout")))
}

func newFilterFields(cfg common.Config) (processors.Processor, error) {
	c := struct {
		Key        string   `config:"key"`
		Topic      string   `config:"topic"`
		Regexp     string   `config:"regexp"`
		Names      []string `config:"names"`
		DateLayout string   `config:"date_layout"`
	}{}
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("fail to unpack the filter_fields configuration:%s", err)
	}

	r, err := regexp.Compile(c.Regexp)
	if err != nil {
		return nil, fmt.Errorf("fail to compile filter_fields regexp:%s error:%v", c.Regexp, err)
	}

	if c.Key == "" {
		c.Key = "message"
	}

	debug("new filter:%v\n", c)

	return filterFields{
		Topic:      c.Topic,
		Key:        c.Key,
		Regexp:     r,
		Names:      c.Names,
		DateLayout: c.DateLayout,
		timeValue:  &TimeValue{},
	}, nil
}

func (f filterFields) Run(event common.MapStr) (common.MapStr, error) {
	val, err := event.GetValue("type")
	if err != nil {
		return event, err
	}

	if f.Topic != val.(string) {
		debug("expect topic:%v, val:%v", f.Topic, val)
		return event, nil
	}

	if val, err = event.GetValue(f.Key); err != nil {
		debug("key:%v not found", f.Key)
		return event, err
	}

	for i, v := range f.Regexp.FindStringSubmatch(fmt.Sprintf("%v", val)) {
		if i == 0 {
			continue
		}
		debug("field[%d]:%v", i, v)
		key := fmt.Sprintf("field_%d", i)
		if i <= len(f.Names) {
			key = f.Names[i-1]
		}
		if key == "@timestamp" && f.DateLayout != "" {
			t, err := time.ParseInLocation(f.DateLayout, v, time.Local)
			if err != nil {
				return event, err
			}
			event.Put(key, f.timeValue.timestamp(t))
			continue
		}

		event.Put(key, v)
	}

	return event, nil
}

func (f filterFields) String() string {
	return fmt.Sprintf("filter_fields=key:%s,regexp:%s,names:%s,date_layout:%s", f.Key, f.Regexp.String(), strings.Join(f.Names, ","), f.DateLayout)
}
