如何使用
典型的应用场景是：想记录一些日志，计时，计数等，可以在web的某个端口方便的查看最近的信息
主要的作用是监控和查错


//创建监控黑板 黑板一共存在三种kit， 分别是logkit counterkit timerkit 用于记录日志，计数器，计时器
bb = blackboardkit.NewBlockBorad(groupname, bbname, bbreadme)


//如何监控
import "github.com/rz1226/rzlib/blackboardkit"
func main(){
    blackboardkit.StartMonitor("9091" ) // 用浏览器看9091/查看监控数据
}



