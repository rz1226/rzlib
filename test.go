package main

import (
	"errors"
	"fmt"
	"github.com/rz1226/rzlib/blackboardkit"
	"reflect"
	"time"
)

type SomeBB struct{
	InsertUser *blackboardkit.BlackBoradKit  `readme:"插入用户信息 "`
	GaoDeApi  *blackboardkit.BlackBoradKit `readme:"调用高德地图api "`
	Db *blackboardkit.BlackBoradKit `readme:"调用数据库 "`

}
//每一个blackboard包含info, err, warn 日志， 一个计数器，一个计时器

func BBinit(dstStruct interface{}, groupName string )(resErr error){
	currentField := ""
	defer func() {
		if co := recover(); co != nil {
			resErr = errors.New("发生panic, field=" + currentField + ":" + fmt.Sprint(co))
		}
	}()

	v := reflect.ValueOf(dstStruct)
	t := v.Type().Elem()
	switch v.Kind() {
	case reflect.Ptr:
		for i := 0; i < v.Elem().NumField(); i++ {
			fieldName := t.Field(i).Name
			tag := t.Field(i).Tag.Get("readme")



			vType := t.Field(i).Type
			fmt.Println("type=", vType )
			fmt.Println(vType  , fieldName, tag )
			if fmt.Sprint(vType)  == "*blackboardkit.BlackBoradKit" {

				bb := blackboardkit.NewBlockBorad(groupName, fieldName, tag)
				v.Elem().Field(i).Set(reflect.ValueOf(bb))
			}else{
				return errors.New("only support int64, float64, string  ")
			}
		}
		return nil
	default:
		return errors.New("only support struct pointer")
	}

}


func main(){

	bb := SomeBB{}

	info := BBinit(&bb, "somegroup" )
	fmt.Println(info )

	blackboardkit.StartMonitor("9090")
	for i:=0;i<5 ;i++  {
		go add( bb )
	}


	time.Sleep( time.Second*1000)

}
func add(bb SomeBB){
	for i:=0;i<10000000000 ;i++  {
		t := bb.Db.Start("开始操作db")
		time.Sleep(time.Microsecond)
		bb.Db.Log("这是db日志", i )
		bb.Db.Inc()
		bb.Db.End(t )


		t2 := bb.InsertUser.Start("开始操作注册")
		time.Sleep(time.Microsecond)
		bb.InsertUser.Log("这是注册用户日志", i )
		bb.InsertUser.Inc()
		bb.InsertUser.End(t2 )
	}
}