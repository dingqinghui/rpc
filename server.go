/**
 * @Author: dingQingHui
 * @Description:
 * @File: server
 * @Version: 1.0.0
 * @Date: 2024/11/19 18:50
 */

package rpc

import (
	"errors"
	"github.com/dingqinghui/extend/component"
	"github.com/dingqinghui/extend/workers"
	"github.com/dingqinghui/zlog"
	"go.uber.org/zap"
	"reflect"
	"unicode"
	"unicode/utf8"
)

type method struct {
	fun       reflect.Value
	typ       reflect.Type
	argType   reflect.Type
	replyType reflect.Type
	haveReply bool
}

func (m *method) call(args []reflect.Value) (err error) {
	function := m.fun
	returnValues := function.Call(args)
	errInter := returnValues[0].Interface()
	if errInter != nil {
		return errInter.(error)
	}
	return nil
}

type service struct {
	name    string
	typ     reflect.Type
	value   reflect.Value
	methods map[string]*method
}

func NewServer(options ...Option) *server {
	s := new(server)
	s.opts = loadOptions(options...)

	return s
}

type server struct {
	component.BuiltinComponent
	opts        *Options
	serviceDict map[string]*service
}

func (s *server) Init() {
	s.serviceDict = make(map[string]*service)
	s.AddComponent(s.opts.GetMsgque())
	s.BuiltinComponent.Init()
}

func (s *server) RegisterName(name string, rcv interface{}) error {
	_, err := s.register(rcv, name)
	s.opts.GetMsgque().Subscribe(name, func(subj string, data []byte) []byte {
		return s.process(subj, data)
	})
	return err
}

func (s *server) process(serName string, data []byte) []byte {
	ser := s.serviceDict[serName]
	if ser == nil {
		zlog.Error("rpc not register service", zap.String("name", serName))
		return nil
	}
	// 解包RPC消息
	message, err := DecodeMessage(data)
	if err != nil {
		zlog.Error("decode rpc message err", zap.Error(err))
		return nil
	}
	md := ser.methods[message.Method]
	if md == nil {
		zlog.Error("service not method", zap.String("service", serName), zap.String("method", message.Method))
		return nil
	}
	// 解码参数
	arg := newByType(md.argType)
	if err := s.opts.GetCodec().Decode(message.Data, arg); err != nil {
		zlog.Error("decode msg err", zap.String("service", serName), zap.String("method", message.Method), zap.Error(err))
		return nil
	}
	haveReply := md.haveReply

	var args []reflect.Value

	if md.argType.Kind() != reflect.Ptr {
		args = append(args, ser.value, reflect.ValueOf(arg).Elem())
	} else {
		args = append(args, ser.value, reflect.ValueOf(arg))

	}
	var reply interface{}
	if haveReply {
		reply = newByType(md.replyType)
		args = append(args, reflect.ValueOf(reply))
	}
	// 执行方法
	workers.Try(func() {
		err = md.call(args)
	}, func(err interface{}) {
		zlog.Error("call err", zap.String("service", serName),
			zap.String("method", message.Method), zap.Error(err.(error)), zap.Stack("stack"))
	})

	if err != nil || !haveReply {
		return nil
	}
	buf, err := s.opts.GetCodec().Encode(reply)
	if err != nil {
		return nil
	}
	return buf
}

func (s *server) register(rcv interface{}, name string) (string, error) {
	ser := new(service)
	ser.typ = reflect.TypeOf(rcv)
	ser.value = reflect.ValueOf(rcv)
	serName := reflect.Indirect(ser.value).Type().Name() // Type

	if serName == "" {
		err := errors.New("server register error name")
		zlog.Error("server register error name", zap.String("name", ser.typ.String()))
		return serName, err
	}
	if !isExported(serName) {
		err := errors.New("server  name is not exported")
		zlog.Error("server register error", zap.Error(err), zap.String("name", ser.typ.String()))
		return serName, err
	}
	ser.name = serName

	// Install the methods
	ser.methods = suitableMethods(serName, ser.typ)
	s.serviceDict[name] = ser
	return serName, nil
}

func (s *server) Stop() {
	s.BuiltinComponent.Stop()
}

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()

func isExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return isExported(t.Name()) || t.PkgPath() == ""
}

func suitableMethods(serName string, typ reflect.Type) map[string]*method {
	methods := make(map[string]*method)
	for m := 0; m < typ.NumMethod(); m++ {
		md := typ.Method(m)
		mtype := md.Type
		mname := md.Name
		// Method must be exported.
		if md.PkgPath != "" {
			continue
		}
		// Method needs four ins: receiver,  *args, *reply.
		if mtype.NumIn() < 2 {
			zlog.Warn("wrong number of input args", zap.String("service", serName),
				zap.String("method", mname), zap.Int("numIn", mtype.NumIn()))
			continue
		}

		argType := mtype.In(1)
		if !isExportedOrBuiltinType(argType) {
			zlog.Warn("parameter type not exported", zap.String("service", serName),
				zap.String("method", mname), zap.String("argType", argType.String()))
			continue
		}

		if mtype.NumOut() != 1 {
			zlog.Error("wrong number of outs", zap.String("service", serName),
				zap.String("method", mname), zap.Int("numOut", mtype.NumOut()))
			continue
		}

		if returnType := mtype.Out(0); returnType != typeOfError {
			zlog.Error("returns not error", zap.String("service", serName),
				zap.String("method", mname), zap.Int("numOut", mtype.NumOut()))
			continue
		}
		mtd := &method{fun: md.Func, argType: argType}
		zlog.Info("", zap.String("method", mname), zap.Int("count", mtype.NumIn()))
		haveReply := mtype.NumIn() > 2
		if haveReply {
			replyType := mtype.In(2)
			if replyType.Kind() != reflect.Ptr {
				zlog.Error("reply type not a pointer", zap.String("service", serName),
					zap.String("method", mname), zap.String("replyType", replyType.String()))
				continue
			}

			if !isExportedOrBuiltinType(replyType) {
				zlog.Error("reply type not exported", zap.String("service", serName),
					zap.String("method", mname), zap.String("replyType", replyType.String()))
				continue
			}
			mtd.replyType = replyType
		}
		mtd.haveReply = haveReply
		methods[mname] = mtd

	}
	return methods
}

func newByType(t reflect.Type) interface{} {
	var argv reflect.Value

	if t.Kind() == reflect.Ptr {
		argv = reflect.New(t.Elem())
	} else {
		argv = reflect.New(t)
	}
	return argv.Interface()
}
