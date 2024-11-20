/**
 * @Author: dingQingHui
 * @Description:
 * @File: message
 * @Version: 1.0.0
 * @Date: 2024/11/20 10:38
 */

package rpc

import "github.com/dingqinghui/extend/codec"

type Message struct {
	Method string
	Data   []byte
}

func EncodeMessage(method string, arg interface{}, bodyCodec codec.ICodec) ([]byte, error) {
	body, err := bodyCodec.Encode(arg)
	if err != nil {
		return nil, err
	}
	m := new(Message)
	m.Method = method
	m.Data = body
	buf, err := codec.Json.Encode(m)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func DecodeMessage(buf []byte) (*Message, error) {
	m := new(Message)
	if err := codec.Json.Decode(buf, m); err != nil {
		return nil, err
	}
	return m, nil
}
