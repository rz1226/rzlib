package ratekit

import (
	"fmt"
	"github.com/rz1226/rzlib/blackboardkit"
	"math/rand"
	"sync/atomic"
	"time"
)

/*
	这个库的功能是以特定的速度异步执行批量闭包函数。
    去掉了retry， 因为没啥用处，反而使得设计更难. 重试是在调用方处理是最明智的。
*/

const RATEKITCHANSIZE = 100               //任务队列chan长度
const MAX_TASK_NUM_EVERY_SECOND = 1000000 //最大速度
const FACTORY_SIZE = 200

type Task struct {
	f func() bool //任务函数，要符合这个格式，实际应用中一般来说是个闭包
}

type RateKit struct {
	tokenChan           chan struct{} //令牌队列
	currentCount        uint32        //上一秒执行的次数
	limitCountPerSecond uint32        //每秒最大执行次数
	asynChan            chan Task
	factory             *WorkerFactory
	bb                  *blackboardkit.BlackBoradKit //记录日志信息用的黑板
}

//limitcount 每秒钟运行次数上限制， workernum工作go程数量， issync是否是同步模式
func NewRateKit(limitCount uint32) *RateKit {
	rk := &RateKit{}
	rk.currentCount = 0
	if limitCount > MAX_TASK_NUM_EVERY_SECOND {
		limitCount = MAX_TASK_NUM_EVERY_SECOND
	}
	rk.limitCountPerSecond = limitCount
	rk.asynChan = make(chan Task, RATEKITCHANSIZE)
	rk.tokenChan = make(chan struct{}, MAX_TASK_NUM_EVERY_SECOND)
	rk.bb = blackboardkit.NewBlockBorad("ratekit", "ratekit", "限速器记录")

	//启动异步任务go程
	rk.factory = newFactory(rk.asynTask, FACTORY_SIZE)
	rk.factory.setRunning()
	go rk.resetAll()
	go rk.releaseToken()

	return rk
}
func (rk *RateKit) incCount() {
	atomic.AddUint32(&rk.currentCount, 1)
}

func (rk *RateKit) resetAll() {
	t1 := time.Tick(time.Second)
	for {
		<-t1
		rk.bb.Log("重置清零 currentcount=" + fmt.Sprint(atomic.LoadUint32(&rk.currentCount)))
		atomic.StoreUint32(&rk.currentCount, 0)
	}

}
func (rk *RateKit) getCount() uint32 {
	return atomic.LoadUint32(&rk.currentCount)
}

func (rk *RateKit) getLimitCount() uint32 {
	return atomic.LoadUint32(&rk.limitCountPerSecond)
}

//放令牌
func (rk *RateKit) releaseToken() {
	limit := rk.getLimitCount()
	t1 := time.Tick(time.Second)
	for {
		<-t1
		now := time.Now()
		for i := 0; i < int(limit); i++ {
			rk.tokenChan <- struct{}{}
		}
		du := time.Since(now)
		rk.bb.Log("释放token=" + fmt.Sprint(limit) + "个，耗时=" + fmt.Sprint(du))
	}
}

//如果chan满了，这个函数会阻塞
//这个函数的基本功能是把函数包装后放入队列
//需要控制入队列的速度，否则队列满的话，会阻塞大批worker在任务未完成重新放入队列的时候
func (rk *RateKit) Go(f func() bool) {
	rk.bb.Inc()

	rf := Task{}
	rf.f = f
	rk.asynChan <- rf
}

//这个就是要送给worker执行任务函数
//这个函数的基本任务是把任务从队列拿出来然后执行
func (rk *RateKit) asynTask() bool {
	rf := <-rk.asynChan
	<-rk.tokenChan //消费令牌
	f := rf.f
	rk.incCount()
	return f()
}

/***************************************worker*******************************/

type WorkerFactory struct {
	workerList []*Worker                    //正在运行的worker集合
	bb         *blackboardkit.BlackBoradKit //记录日志信息用的黑板

}

//第一个参数是worker工作的函数
func newFactory(f func() bool, size int) *WorkerFactory {
	fac := &WorkerFactory{}
	fac.workerList = make([]*Worker, 0, size)
	for i := 0; i < size; i++ {
		worker := &Worker{}
		worker.f = f
		worker.factory = fac
		worker.setStop()
		fac.workerList = append(fac.workerList, worker)
	}

	fac.bb = blackboardkit.NewBlockBorad("ratekit", "workerfactory", "worker工厂记录")

	return fac
}

func (fac *WorkerFactory) setRunning() {
	for _, v := range fac.workerList {
		v.setRunning()
	}
}

type Worker struct {
	running uint32      //0:关闭  1：打开
	f       func() bool //工作函数，例如从某个chan拿出数据然后处理，可以一直重复执行 ,函数结束后协程退出
	//数据并且处理
	factory *WorkerFactory //所属的工厂
}

func (w *Worker) isRunning() bool {
	v := atomic.LoadUint32(&w.running)
	if v == 0 {
		return false
	}
	return true
}
func (w *Worker) setRunning() {
	atomic.StoreUint32(&w.running, 1)
	go w.run()
}
func (w *Worker) setStop() {
	atomic.StoreUint32(&w.running, 0)
}

func (w *Worker) run() {
	defer func() {
		if co := recover(); co != nil {
			w.factory.bb.Panic("worker_panic", "worker 发生异常:", co)
			time.Sleep(time.Millisecond * time.Duration(20+getRandInt(300))) //如果挂了，等一点点时间再重启，防止无限挂跑死cpu
			w.run()
		}
	}()

	if w.isRunning() == true {
		for {
			result := w.f()
			if result == true {
				w.factory.bb.Log("worker 函数执行成功")
			} else {
				w.factory.bb.Err("worker 函数执行失败")
			}
			if w.isRunning() == false {
				//状态为停止，关闭worker
				w.factory.bb.Warn("关闭 worker")
				return
			}
		}
	}

}

func getRandInt(n int) int {
	rand.Seed(time.Now().Unix())
	r := rand.Intn(n)
	return r
}
