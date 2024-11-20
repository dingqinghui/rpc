/**
 * @Author: dingQingHui
 * @Description:
 * @File: Con
 * @Version: 1.0.0
 * @Date: 2024/11/19 10:15
 */

package nats

import (
	"github.com/dingqinghui/extend/component"
	"github.com/dingqinghui/extend/workers"
	"github.com/dingqinghui/rpc/common"
	"github.com/dingqinghui/zlog"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"time"
)

func New(options ...Option) *Con {
	c := new(Con)
	c.opts = loadOptions(options...)
	c.msgChan = make(chan *nats.Msg, 4096)
	return c
}

type Con struct {
	component.BuiltinComponent
	opts    *Options
	rawCon  *nats.Conn
	msgChan chan *nats.Msg
}

func (c *Con) Name() string {
	return "nats"
}

func (c *Con) Init() {
	c.connect()
}

func (c *Con) Stop() {
	c.rawCon.Close()
	close(c.msgChan)
}

func (c *Con) connect() {
	url := c.opts.GetUrl()
	con, err := nats.Connect(url, c.opts.natsOptions...)
	if err != nil {
		zlog.Error("nats connect err", zap.String("url", url), zap.Error(err))
		return
	}
	c.rawCon = con
	zlog.Info("nats connect", zap.String("url", url), zap.Error(err))
}

// Call
// @Description: 同步发送请求
// @receiver c
// @param subj
// @param data
// @param timeout
// @return *nats.Msg
// @return error
func (c *Con) Call(subj string, data []byte, timeout time.Duration) ([]byte, error) {
	msg, err := c.rawCon.Request(subj, data, timeout)
	if err != nil {
		return nil, err
	}
	return msg.Data, err
}

// Send
// @Description: 异步调用
// @receiver c
// @param subj
// @param data
// @return err
func (c *Con) Send(subj string, data []byte) (err error) {
	if err = c.rawCon.Publish(subj, data); err != nil {
		zlog.Error("nats publish error", zap.Error(err))
		return
	}
	return
}

func (c *Con) Subscribe(subject string, process common.ProcessFunc) {
	_, chanErr := c.rawCon.ChanSubscribe(subject, c.msgChan)
	if chanErr != nil {
		zlog.Error("nats chan subscribe error", zap.Error(chanErr))
		return
	}

	workers.Submit(func() {
		for msg := range c.msgChan {
			reply := process(msg.Subject, msg.Data)
			if reply == nil || msg.Reply == "" {
				continue
			}
			if err := msg.Respond(reply); err != nil {
				zlog.Error("nats chan respond error", zap.Error(err))
				continue
			}
		}

	}, func(err interface{}) {
		zlog.Panic("nats process panic", zap.Error(err.(error)))
	})
}
