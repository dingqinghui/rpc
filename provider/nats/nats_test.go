/**
 * @Author: dingQingHui
 * @Description:
 * @File: nats_test
 * @Version: 1.0.0
 * @Date: 2024/11/19 10:49
 */

package nats

import (
	"github.com/dingqinghui/zlog"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestNatsPublish(t *testing.T) {
	c := New(WithUrl("127.0.0.1:4222"))
	c.Init()
	c.Subscribe("top1", func(msg *nats.Msg) {
		zlog.Info("msg", zap.String("data", string(msg.Data)))
	})

	c.Subscribe("top1", func(msg *nats.Msg) {
		zlog.Info("msg", zap.String("data", string(msg.Data)))
	})

	c.ChanExecute("top2", func(msg *nats.Msg) {
		zlog.Info("msg", zap.String("data", string(msg.Data)))
	})

	c.Send("top1", []byte("11111111"))

	for i := 0; i < 3; i++ {
		c.Send("top2", []byte("22222222222222"))
	}

	time.Sleep(time.Second * 10)
}
