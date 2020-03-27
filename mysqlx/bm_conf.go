package mysqlx

import "reflect"

//一些全局配置
var Conf Config

type Config struct {
	TagName string //业务模型映射到数据库字段的tag名,影响查询结果转数据模型，以及数据模型转Line Lines

	//判断是不是auto属性，也就是字段值是数据库自动生成的，例如auto_increment , create_time ,update_time等 影响插入
	//第一个参数包含了所有的tag ,可以通过 tags.Get(name)获取某个tag
	//影响业务模型数据转化为Line  Lines
	FuncAuto func(tags reflect.StructTag) bool
}

func init() {
	Conf = Config{}
	//默认的数据库映射tag = orm
	Conf.TagName = "orm"
	//判断是不是auto字段的默认规则是  `auto:“1”`
	f := func(tags reflect.StructTag) bool {
		tag := tags.Get("auto")
		if tag == "1" {
			return true
		}
		return false
	}
	Conf.FuncAuto = f
}
