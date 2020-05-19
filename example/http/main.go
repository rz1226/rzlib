package main

import (
	"github.com/rz1226/rzlib/blackboardkit"
	"github.com/rz1226/rzlib/httpkit"
	"time"
)

var client *httpkit.HTTPClient

func init() {

	client = httpkit.NewHTTPClient(1, 100)
	blackboardkit.StartMonitor("9090")
}

func main() {
	for i := 0; i < 100; i++ {
		client.Get("http:// www.baidu.com")
	}

	time.Sleep(time.Second * 1000)

}
