package errset

import (
	"errors"
	"fmt"
)

//  错误集合   可以记录多个错误，  例如错误累加的场景
const NEWLINE = "\n"

type ErrSet struct {
	errs []error
	logs []string
}

func NewErrSet() *ErrSet {
	es := new(ErrSet)
	es.errs = make([]error, 0, 10)
	es.logs = make([]string, 0, 10)
	return es
}

func (es *ErrSet) Log(str string) *ErrSet {
	es.logs = append(es.logs, str)
	return es
}

func (es *ErrSet) Add(err error) *ErrSet {
	es.errs = append(es.errs, err)
	return es
}

func (es *ErrSet) HaveErr() bool {
	return es.Len() != 0
}

func (es *ErrSet) Len() int {
	return len(es.errs)
}

func (es *ErrSet) Error() error {
	if es.HaveErr() {
		str := es.ErrStr()
		return errors.New(str)
	}
	return nil
}

func (es *ErrSet) ErrStr() string {
	str := "[ErrSet:"
	for k, v := range es.errs {
		str += fmt.Sprint("第", k, "个错误：", v.Error())
		str += NEWLINE
	}

	for k, v := range es.logs {
		str += fmt.Sprint("第", k, "个日志：", v)
		str += NEWLINE
	}
	str += "ErrSet]"
	return str
}
