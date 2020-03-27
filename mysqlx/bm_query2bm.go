package mysqlx

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

//多个，把queryRes里的数据转化成*[]*struct
func (q *QueryRes) ToStruct(dstStructs interface{}) error {
	if q.err != nil {
		return q.err
	}
	data := q.Data()
	if len(data) == 0 {
		return errors.New("empty")
	}
	isMany, err := isBmMany(dstStructs)
	if err != nil {
		return err
	}
	if isMany == true {
		return queryRes2StructBatch(data, dstStructs, nil)
	} else {
		//处理单个
		return queryRes2Struct(data[0], dstStructs, nil)
	}

}

// map数据映射到struct, 同名映射到struct的key不区分大小写，检查类型  。dstStruct是接收数据的struct的指针
//可以加多个struct接收数据, 一般只支持 整数，浮点，字符串
// 第二个参数是struct 指针
func queryRes2Struct(sourceData map[string]interface{}, dstStruct interface{}, f func(map[string]interface{})) error {
	if f != nil {
		f(sourceData)
	}
	err := structFromQueryRes(sourceData, dstStruct)
	if err != nil {
		return err
	}
	return nil
}

//第二个参数是&[]*SomeStruct
func queryRes2StructBatch(sourceDatas []map[string]interface{}, dstStructs interface{}, f func(map[string]interface{})) error {
	strusRV := reflect.Indirect(reflect.ValueOf(dstStructs))
	elemRT := strusRV.Type().Elem()
	for _, v := range sourceDatas {
		eleData := reflect.New(elemRT.Elem()).Interface()
		err := queryRes2Struct(v, eleData, f)
		if err != nil {
			return err
		}
		strusRV = reflect.Append(strusRV, reflect.ValueOf(eleData))
	}
	reflect.Indirect(reflect.ValueOf(dstStructs)).Set(strusRV)
	return nil
}

//only support int64, float64, string, []byte
func structFromQueryRes(sourceData map[string]interface{}, dstStruct interface{}) (resErr error) {
	//当前处理到哪个key了。panic返回报错用的
	currentField := ""
	defer func() {
		if co := recover(); co != nil {
			resErr = errors.New("发生panic, field=" + currentField + ":" + fmt.Sprint(co))
		}
	}()
	length := len(sourceData)
	if length <= 0 {
		return errors.New("no sourceData ,len zero ")
	}
	v := reflect.ValueOf(dstStruct)
	t := v.Type().Elem()
	switch v.Kind() {
	case reflect.Ptr:
		for i := 0; i < v.Elem().NumField(); i++ {
			//oriKey := t.Field(i).Name
			key := t.Field(i).Tag.Get(Conf.TagName)

			if len(key) == 0 {
				//找不到业务模型struct的数据库映射tag,忽略
				continue
			}
			currentField = key
			valueFromMap, ok := sourceData[key]
			if !ok {
				continue
			}
			vType := t.Field(i).Type
			switch vType.Name() {
			case "int64":
				valueInt64, ok := valueFromMap.(int64)

				if !ok {
					if valueStr, ok := valueFromMap.(string); ok {
						valueInt64, err := strconv.ParseInt(valueStr, 10, 64)

						if err == nil {

							v.Elem().Field(i).Set(reflect.ValueOf(valueInt64))
						}
					} else {
						return errors.New("field " + key + " can not store as integer , is " + fmt.Sprint(reflect.TypeOf(valueFromMap)))
					}

				} else {
					v.Elem().Field(i).Set(reflect.ValueOf(valueInt64))
				}

			case "float64":
				valueF64, ok := valueFromMap.(float64)
				if !ok {
					if valueStr, ok := valueFromMap.(string); ok {
						valueInt64, err := strconv.ParseFloat(valueStr, 64)
						if err == nil {
							v.Elem().Field(i).Set(reflect.ValueOf(valueInt64))
						}
					} else {
						return errors.New("field " + key + " can not store as float ,is " + fmt.Sprint(reflect.TypeOf(valueFromMap)))
					}

				} else {
					v.Elem().Field(i).Set(reflect.ValueOf(valueF64))
				}

			case "string":
				valueString, ok := valueFromMap.(string)
				if !ok {
					//如果不是string类型，就强制转化
					valueString = fmt.Sprint(valueFromMap)
				}
				v.Elem().Field(i).Set(reflect.ValueOf(valueString))
			default:
				return errors.New("only support int64, float64, string  ")
			}
		}
		return nil
	default:
		return errors.New("only support struct pointer")
	}
}
