/**
 * @Author: dingQingHui
 * @Description:
 * @File: api
 * @Version: 1.0.0
 * @Date: 2024/11/20 15:05
 */

package common

import (
	"github.com/dingqinghui/extend/component"
	"time"
)

type ProcessFunc func(subj string, data []byte) []byte
type IMessageQue interface {
	component.IComponent
	Call(subj string, data []byte, timeout time.Duration) ([]byte, error)
	Send(subj string, data []byte) (err error)
	Subscribe(subject string, process ProcessFunc)
}
