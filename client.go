/**
 * @Author: dingQingHui
 * @Description:
 * @File: Client
 * @Version: 1.0.0
 * @Date: 2024/11/19 17:57
 */

package rpc

import (
	"github.com/dingqinghui/extend/component"
	"time"
)

func NewClient(options ...Option) *Client {
	c := new(Client)
	c.opts = loadOptions(options...)
	return c
}

type Client struct {
	component.BuiltinComponent
	opts *Options
}

func (c *Client) Init() {
	c.AddComponent(c.opts.GetMsgque())
	c.BuiltinComponent.Init()
}

func (c *Client) Send(service string, method string, msg interface{}) error {
	buf, err := EncodeMessage(method, msg, c.opts.GetCodec())
	if err != nil {
		return err
	}
	return c.opts.GetMsgque().Send(service, buf)
}

func (c *Client) Call(service string, method string, msg interface{}, reply interface{}, timeout time.Duration) error {
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

func (c *Client) Stop() {
	c.BuiltinComponent.Stop()
}
