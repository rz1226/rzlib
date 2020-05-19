httpclient用法
典型的应用场景，访问一个http url

var client *httpkit.HttpClient

func init(){
    1 超时时间秒 2 最大空闲连接池容量
	client = httpkit.NewHttpClient(1,100)
}
以上初始化一个全局的httpclient

// 见example/httpkit

// 启动监控
blackboardkit.StartMonitor("9091" ) 


