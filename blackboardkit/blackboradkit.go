package blackboardkit

import (
	"fmt"
	"github.com/rz1226/rzlib/kits"
	"github.com/rz1226/rzlib/serverkit"
	"net/http"
	"sync"
	"time"
)

/**
理想的api是
type SomeBB struct{
	InsertUser BlackBoard `readme:"插入用户信息日志"`
	GaoDeApi  BlackBoard `readme:"调用高德地图api的日志"`
	Db BlackBoard `readme:"调用数据库错误"`

}
每一个blackboard包含info, err, warn 日志， 一个计数器，一个计时器

bb := SomeBB{}
BBinit(&bb, groupName )

type


bb.Info.Info("xx)  .Err()  .Warn()
bb.Info.Inc()  .IncBy(1)
t := bb.Info.Start()
bb.Info.Ends(t)

有类型保护, 调用的时候不容易出错
降低了使用成本

*/

var allbb *AllBB

func init() {
	allbb = sNewAllBB()

}
func ShowAllBBs() string {
	return allbb.showAll()
}
func ShowGroup(groupName string) string {
	return allbb.show(groupName)
}

//所有的bb
type AllBB struct {
	data map[string][]*BlackBoradKit
	mu   *sync.Mutex
}

func sNewAllBB() *AllBB {
	a := &AllBB{}
	a.mu = &sync.Mutex{}
	a.data = make(map[string][]*BlackBoradKit, 0)
	return a
}
func (a *AllBB) add(bb *BlackBoradKit) {
	a.mu.Lock()
	defer a.mu.Unlock()
	groupName := bb.groupname
	_, ok := a.data[groupName]
	if !ok {
		a.data[groupName] = make([]*BlackBoradKit, 0, 5)
	}
	a.data[groupName] = append(a.data[groupName], bb)
}
func (a *AllBB) show(groupName string) string {
	a.mu.Lock()
	defer a.mu.Unlock()
	str := ""
	v, ok := a.data[groupName]
	if !ok {
		return "没找到此分组:" + groupName
	}
	for _, v2 := range v {
		str += v2.show()
	}

	return str
}
func (a *AllBB) showAll() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	str := ""
	for _, v := range a.data {
		for _, v2 := range v {
			str += v2.show()
		}

	}
	return str
}

//监控信息黑板

type BlackBoradKit struct {
	logKit           *kits.LogKit
	timerKit         *kits.TimerKit
	counterKit       *kits.CounterKit
	bbStartTime      string
	noPrintToConsole bool
	name             string // name readme 传递到底层
	readme           string
	groupname        string
}

func NewBlockBorad(groupName, name, readme string) *BlackBoradKit {
	bb := &BlackBoradKit{}
	bb.groupname = groupName
	bb.name = name
	bb.readme = readme
	bb.bbStartTime = time.Now().Format("2006-01-02 15:04:05")
	bb.initLogKit()
	bb.initCounterKit()
	bb.initTimerKit()
	bb.noPrintToConsole = true //默认不直接打印信息
	allbb.add(bb)
	return bb
}

//是否同时打印到标准输出
func (bb *BlackBoradKit) SetNoPrintToConsole(result bool) {
	bb.noPrintToConsole = result
}

//初始化日志kit
func (bb *BlackBoradKit) initLogKit() {
	bb.logKit = kits.NewLogKit(bb.name, bb.readme)
}

//初始化计数器kit
func (bb *BlackBoradKit) initCounterKit() {
	bb.counterKit = kits.NewCounterKit(bb.name, bb.readme)
}

//初始化计时器kit
func (bb *BlackBoradKit) initTimerKit() {
	bb.timerKit = kits.NewTimerKit(bb.name, bb.readme)
}

/*----------------------------log--------------------------------*/
func (bb *BlackBoradKit) Log(logs ...interface{}) {
	str := bb.logKit.PutContentsAndFormat(logs...)
	if bb.noPrintToConsole == false {
		fmt.Print(str)
	}
}

/*---------------------------timer---------------------------------*/
func (bb *BlackBoradKit) Start(tickInfo string) *kits.Tick {
	return bb.timerKit.Start(tickInfo)
}
func (bb *BlackBoradKit) End(tick *kits.Tick) {
	bb.timerKit.End(tick)
}

/*---------------------------counter---------------------------------*/
func (bb *BlackBoradKit) Inc() {
	bb.counterKit.Inc()
}
func (bb *BlackBoradKit) IncBy(num int64) {
	bb.counterKit.IncBy(num)
}

/*--------------------------show---------------------------*/
//获取监控信息
func (bb *BlackBoradKit) show() string {

	str := "\n\n\n----------------" + bb.name + " blackboard info ----------------- : \n\n\n"

	str += "监控启动时间:" + bb.bbStartTime + "\n"
	str += bb.logKit.Show()

	str += "\n\n\n"
	str += bb.counterKit.Show()

	str += "\n\n\n"
	str += bb.timerKit.Show()
	return str
}

func httpShowAll(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("yes")
	r.ParseForm()
	str := ShowAllBBs()
	fmt.Fprintln(w, str)
}

func StartMonitor(port string) {
	go serverkit.NewSimpleHttpServer().Add("/", httpShowAll).Start(port)
}
