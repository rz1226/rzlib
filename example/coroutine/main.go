package main

import (
	"fmt"
	"github.com/rz1226/rzlib/coroutinekit"
	"time"
)

func test() {
	for {
		time.Sleep(time.Millisecond)
		fmt.Println("xxx")
		panic(1)
	}
}

func main() {
	coroutinekit.StartMonitor("9090")
	coroutinekit.Start("测试", 2, test, true)

	fmt.Println("done")
	time.Sleep(time.Second * 1000)
}
