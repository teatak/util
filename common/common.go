package common

import (
	"context"
	"math/rand"
	"time"
)

const MobileJoiner = ":"
const WechatUnionPrefix = "wechat:"
const Date = "2006-01-02"
const DateTime = "2006-01-02 15:04:05"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Context(dur ...time.Duration) (context.Context, context.CancelFunc) {
	var duration time.Duration = 5 * time.Second
	if len(dur) > 0 {
		duration = dur[0]
	}
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	return ctx, cancel
}
