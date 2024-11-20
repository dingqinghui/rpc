/**
 * @Author: dingQingHui
 * @Description:
 * @File: client
 * @Version: 1.0.0
 * @Date: 2024/11/19 17:57
 */

package rpc

import (
	"github.com/dingqinghui/extend/component"
	"time"
)

func NewClient(options ...Option) *client {
	c := new(client)
	c.opts = loadOptions(options...)
	c.Init()
	return c
}

type client struct {
	component.BuiltinComponent
	opts *Options
}

func (c *client) Init() {
	c.AddComponent(c.opts.GetMsgque())
	c.BuiltinComponent.Init()
}

func (c *client) Send(service string, method string, msg interface{}) error {
	buf, err := EncodeMessage(method, msg, c.opts.GetCodec())
	if err != nil {
		return err
	}
	return c.opts.GetMsgque().Send(service, buf)
}

func (c *client) Call(service string, method string, msg interface{}, reply interface{}, timeout time.Duration) error {
	buf, err := EncodeMessage(method, msg, c.opts.GetCodec())
	if err != nil {
		return err
	}
	replyBuf, err := c.opts.GetMsgque().Call(service, buf, timeout)
	if err != nil {
		return err
	}
	if err := c.opts.GetCodec().Decode(replyBuf, reply); err != nil {
		return err
	}
	return nil
}

func (c *client) Stop() {
	c.BuiltinComponent.Stop()
}
