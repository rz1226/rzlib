package actor

//
//import (
//
//	"fmt"
//	"github.com/rz1226/simplegokit/coroutinekit"
//	"sync/atomic"
//
//	"work/utilx/actor"
//	"github.com/rz1226/simplegokit/blackboardkit"
//	"time"
//)
//func init() {
//
//	coroutinekit.StartMonitor("9090")
//	blackboardkit.StartMonitor("9091")
//
//}
//func main(){
//	f := func(data interface{}) (interface{},error){
//		s := fmt.Sprint(data ) + "_"
//		return s, nil
//	}
//
//
//
//	f5 := func(data interface{}) (interface{},error){
//		fmt.Println("--")
//		fmt.Println("len===" ,len(data.([]interface{}))  )
//
//		fmt.Println("data==",data)
//		return data , nil
//	}
//
//	var  count uint64 = 0
//	f6 := func(data interface{}) (interface{},error){
//		tmp := data.([]interface{})
//
//		atomic.AddUint64(&count,uint64(len(tmp )))
//
//		fmt.Println("一共是--", atomic.LoadUint64(&count ))
//
//		return nil, nil
//
//
//	}
//
//	a :=  actor.NewActor(nil ,10 ,"初始"  ) ;
//	a.AddActor(f,10, "变成数组" ).SetCumulateCount(100).
//		AddActor(f5, 10, "显示" ).
//		AddActor(f6, 1, "数一数对不对"  )
//	a.Run() ;
//	for i:=1;i<100000;i++{
//
//		a.Put(i)
//
//	}
//
//
//	time.Sleep(time.Second*100)
//}
//
