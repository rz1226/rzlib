package actor

import (
	"errors"
	"fmt"
	"github.com/rz1226/rzlib/blackboardkit"
	"github.com/rz1226/rzlib/coroutinekit"
	"sync/atomic"
	"time"
)

var currentId uint64 = 0
var actorbb *blackboardkit.BlackBoradKit
var CONFIG_CSize int //数据队列的长度
var CONFIG_CUM_TICK_LONG time.Duration

func init() {
	CONFIG_CSize = 10
	CONFIG_CUM_TICK_LONG = time.Millisecond * 50
}

func init() {
	actorbb = blackboardkit.NewBlockBorad("actor", "actor", "actor记录")
}

//模拟的actor
type Actor struct {
	Id              uint64
	Name            string
	CumulateCount   uint32
	C               chan interface{}
	F               func(interface{}) (interface{}, error)
	NumOfConcurrent uint8 //并发数量
	Next            []*Actor
}

func NewActor(f func(interface{}) (interface{}, error), numOfConcurrent int, name string) *Actor {
	a := &Actor{}
	a.Name = name

	a.CumulateCount = 1
	a.Id = atomic.AddUint64(&currentId, 1)
	a.C = make(chan interface{}, CONFIG_CSize)
	a.F = f
	a.NumOfConcurrent = uint8(numOfConcurrent)
	if a.NumOfConcurrent <= 0 {
		a.NumOfConcurrent = 1
	}
	a.Next = make([]*Actor, 0, 10)
	return a
}

//不设置，默认为1 表示不累积
func (a *Actor) SetCumulateCount(cumulateCount uint32) *Actor {
	if cumulateCount == 0 {
		cumulateCount = 1
	}
	if cumulateCount > 10000 {
		cumulateCount = 10000
	}
	a.CumulateCount = cumulateCount
	return a
}

func (a *Actor) AddActor(f func(interface{}) (interface{}, error), numOfConcurrent int, name string) *Actor {
	b := NewActor(f, numOfConcurrent, name)
	a.setNext(b)
	return b
}

func (a *Actor) setNext(b *Actor) {
	a.Next = append(a.Next, b)
}

func (a *Actor) Run() {
	a.run()
}
func (a *Actor) Put(data interface{}) {
	a.C <- data
}
func (a *Actor) PutWait(data interface{}, wait int) error {
	timer := time.After(time.Second * time.Duration(wait))
	select {
	case a.C <- data:
		return nil
	case <-timer:
		return errors.New("failed to put in data to actor")
	}
}

func (a *Actor) run() {
	workF := func() {
		var cumulateData []interface{}
		var needCum bool = false
		cumcount := uint32(0)
		if a.CumulateCount > 1 {
			cumulateData = make([]interface{}, 0, a.CumulateCount)
			needCum = true
		}
		t := time.NewTicker(CONFIG_CUM_TICK_LONG)
		for {
			//从队列获取数据
			select {
			case data := <-a.C:
				if data == nil {
					continue
				}
				var res interface{}
				var err error
				if a.F != nil {
					actorbb.Inc()
					t := actorbb.Start("actor_function_run")

					res, err = a.F(data)
					actorbb.End(t)
					if err != nil {
						actorbb.Err("actor_error", err)
						continue
					}
				} else {
					res = data
				}
				if needCum == true {
					cumulateData = append(cumulateData, res)
					cumcount++
					if cumcount == a.CumulateCount {
						for _, v := range a.Next {
							v.Put(cumulateData)
						}
						cumulateData = make([]interface{}, 0, a.CumulateCount)
						cumcount = uint32(0)
					}
				} else {
					for _, v := range a.Next {
						v.Put(res)
					}
				}
			case <-t.C:
				if cumcount == 0 {
					continue
				}
				for _, v := range a.Next {
					v.Put(cumulateData)
				}
				cumulateData = make([]interface{}, 0, a.CumulateCount)
				cumcount = uint32(0)
			}
		}
	}
	coroutinekit.Start("actor job name = "+a.Name+" id="+fmt.Sprint(a.Id), int(a.NumOfConcurrent), workF, true)
	if len(a.Next) > 0 {
		for _, v := range a.Next {
			v.run()
		}
	}
}
