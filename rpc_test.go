/**
 * @Author: dingQingHui
 * @Description:
 * @File: rpc_test
 * @Version: 1.0.0
 * @Date: 2024/11/19 18:05
 */

package rpc

import (
	"fmt"
	"github.com/dingqinghui/extend/codec"
	"testing"
	"time"
)

type TestMessageReq struct {
	A int
}

func TestClient(t *testing.T) {
	var a = 1
	s, _ := codec.Json.Encode(&a)
	println(s)
	c := NewClient()
	c.Init()
	c.Send("TestService", "Add", &TestMessageReq{A: 2})
	res := &TestMessageReq{}
	c.Call("TestService", "Sub", &TestMessageReq{A: 100}, res, time.Second)
	fmt.Printf("res:%v\n", res.A)
}

type TestService struct {
}

func (t *TestService) Add(req *TestMessageReq) error {
	req.A += 1
	fmt.Printf("Add req:%v\n", req.A)
	return nil
}

func (t *TestService) Sub(req *TestMessageReq, res *TestMessageReq) error {
	res.A = req.A - 1
	fmt.Printf("Sub res:%v\n", res.A)
	return nil
}

func TestSerer(t *testing.T) {
	s := NewServer()
	s.Init()
	s.RegisterName("TestService", &TestService{})
	time.Sleep(time.Hour)
}
