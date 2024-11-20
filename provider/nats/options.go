/**
 * @Author: dingQingHui
 * @Description:
 * @File: option
 * @Version: 1.0.0
 * @Date: 2024/11/19 10:17
 */

package nats

import "github.com/nats-io/nats.go"

type Option func(*Options)

func loadOptions(options ...Option) *Options {
	opts := new(Options)
	for _, option := range options {
		option(opts)
	}
	return opts
}

type Options struct {
	url         string
	natsOptions []nats.Option
}

func (o *Options) GetUrl() string {
	if o.url != "" {
		return o.url
	}
	return "127.0.0.1:4222"
}

func WithUrl(url string) Option {
	return func(op *Options) {
		op.url = url
	}
}
