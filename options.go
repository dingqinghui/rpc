/**
 * @Author: dingQingHui
 * @Description:
 * @File: options
 * @Version: 1.0.0
 * @Date: 2024/11/20 15:10
 */

package rpc

import (
	"github.com/dingqinghui/extend/codec"
	"github.com/dingqinghui/rpc/common"
	"github.com/dingqinghui/rpc/provider/nats"
)

type Option func(*Options)

func loadOptions(options ...Option) *Options {
	opts := new(Options)
	for _, option := range options {
		option(opts)
	}
	return opts
}

type Options struct {
	msgque common.IMessageQue
	codec  codec.ICodec
}

func (o *Options) GetMsgque() common.IMessageQue {
	if o.msgque == nil {
		o.msgque = nats.New()
	}
	return o.msgque
}
func (o *Options) GetCodec() codec.ICodec {
	if o.codec != nil {
		return o.codec
	}
	return codec.Json
}

func WithMessageQue(msgque common.IMessageQue) Option {
	return func(op *Options) {
		op.msgque = msgque
	}
}

func WithCodec(codec codec.ICodec) Option {
	return func(op *Options) {
		op.codec = codec
	}
}
