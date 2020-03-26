package main

import (
	"fmt"
	"github.com/rz1226/rzlib/blackboardkit"
	"github.com/rz1226/rzlib/ratekit"
	rawhttp "net/http"
	_ "net/http/pprof"
	"time"
)

func init() {

	go func() {
		rawhttp.ListenAndServe(":6060", nil)
	}()

}

var rk *ratekit.RateKit

func init() {
	rk = ratekit.NewRateKit(6000 * 1000)

}

func testx() {

	for i := 0; i < 10000000000; i++ {
		f := func(a int) func() bool {
			return func() bool {

				//fmt.Println("闭包", a   )
				//time.Sleep(time.Millisecond*1)
				//panic(1)
				return true
			}

		}(i)

		//time.Sleep(time.Microsecond*20)
		rk.Go(f)

	}
	//fmt.Println( rk.Show())
}

func main() {
	blackboardkit.StartMonitor("9090")
	fmt.Print("haha")
	go func() {

		testx()

	}()

	time.Sleep(time.Second * 300)

}
