/**
 * @Author: dingQingHui
 * @Description:
 * @File: reflect_test
 * @Version: 1.0.0
 * @Date: 2024/11/19 15:24
 */

package nats

import (
	"fmt"
	"reflect"
	"testing"
)

type User struct {
	name string
	age  int
	addr string
}

func TestReflect(t *testing.T) {

	uu := User{"tom", 27, "beijing"}
	u := &uu
	v := reflect.ValueOf(u)

	_t := reflect.TypeOf(u)
	fmt.Println("TypeOf=", _t)

	t1 := reflect.Indirect(v).Type()
	fmt.Println("t1=", t1)

	fmt.Println("t1=", _t.Elem())

	fmt.Println("t1=", _t.PkgPath())

}
