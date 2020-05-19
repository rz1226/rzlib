package ejson

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Ejson struct {
	j      interface{} //  json对象，一般是[]interface{}或者map[string]interface{}
	s      []byte      //  json字符串
	existJ bool
	existS bool
	jErr   bool
	sErr   bool
}

func (e *Ejson) getj() interface{} {
	err := e.setj()

	if err != nil {
		return nil
	}
	return e.j
}
func (e *Ejson) gets() []byte {
	err := e.sets()
	if err != nil {
		return []byte("")
	}
	return e.s
}

func (e *Ejson) setj() error {
	if e.existJ {
		return nil
	}
	if e.jErr {
		return nil
	}
	if !e.existS {
		return errors.New("ejson setj error")
	}
	err := json.Unmarshal(e.s, &e.j)
	if err != nil {
		e.jErr = true
		return err
	}
	e.existJ = true
	return nil
}
func (e *Ejson) sets() error {
	if e.existS {
		return nil
	}
	if e.sErr {
		return nil
	}
	if !e.existJ {
		return errors.New("ejson sets error")
	}
	s, err := json.Marshal(e.j)
	if err != nil {
		e.sErr = true
		return err
	}
	e.s = s
	e.existS = true
	return nil
}

//  和newejson的区别是会检查错误，性能稍低
func NewEjsonWithCheck(j interface{}) (*Ejson, error) {
	json2, err := NewEjson(j)
	if err != nil {
		return nil, err
	}
	if err := json2.setj(); err != nil {
		return nil, err
	}
	if err := json2.sets(); err != nil {
		return nil, err
	}
	return json2, nil
}

func NewEjson(j interface{}) (*Ejson, error) {
	e := Ejson{}
	e.existJ = false
	e.existS = false
	e.jErr = false
	e.sErr = false
	if jdata, ok := j.(string); ok {
		e.s = []byte(jdata)
		e.existS = true
		return &e, nil
	}
	if jdata, ok := j.([]byte); ok {
		e.s = jdata
		e.existS = true
		return &e, nil
	}
	_, ok := j.([]interface{})
	_, ok2 := j.(map[string]interface{})
	if (!ok) && (!ok2) {
		return nil, errors.New("not slice or map, invalid json")
	}
	e.j = j
	e.existJ = true
	return &e, nil
}

// func (e *Ejson) Fresh() {
// 	s, _ := json.Marshal(e.j)
// 	e.s = string(s)
// }

func (e *Ejson) JSON() interface{} {
	return e.getj()
}

func (e *Ejson) Data() interface{} {
	return e.getj()

}
func (e *Ejson) Bytes() []byte {
	return e.gets()
}
func (e *Ejson) String() string {
	return string(e.gets())
}

//  拷贝map类型的json , 就是新创建，用s解析填充
func (e *Ejson) CopyMap() (map[string]interface{}, error) {
	if e.IsMap() {
		var result map[string]interface{}
		err := json.Unmarshal(e.gets(), &result)
		if err != nil {
			return nil, errors.New("ejson:copymap: there is error :" + err.Error())
		}
		return result, nil
	}
	return nil, errors.New("ejson:copymap: this is not map")
}

//  如果不是map类型返回空值
func (e *Ejson) MapData() map[string]interface{} {
	data, ok := e.getj().(map[string]interface{})
	if ok {
		return data
	}
	return make(map[string]interface{})
}

//  如果不是数组类型返回空值
func (e *Ejson) ArrayData() []interface{} {
	data, ok := e.getj().([]interface{})
	if ok {
		return data
	}
	return make([]interface{}, 0)
}

func (e *Ejson) Merge(a *Ejson) (*Ejson, error) {
	if !e.IsMap() || !a.IsMap() {
		return nil, errors.New("ejson:merge 两个值必须都是map")
	}
	obj := &Ejson{j: make(map[string]interface{})}
	for k, v := range e.getj().(map[string]interface{}) {
		obj.j.(map[string]interface{})[k] = v
	}
	for k2, v2 := range a.getj().(map[string]interface{}) {
		obj.j.(map[string]interface{})[k2] = v2
	}
	return obj, nil
}

func (e *Ejson) Len() int {
	if e.IsMap() {
		return len(e.getj().(map[string]interface{}))
	}
	if e.IsArray() {
		return len(e.getj().([]interface{}))
	}
	return 0
}

func (e *Ejson) ArrayGetData(key int) (*Ejson, error) {

	v, ok := e.getj().([]interface{})
	if !ok {
		return nil, errors.New("ejson:ArrayGetData主体对象必须是array")
	}
	if key > (len(v) - 1) {
		return nil, errors.New("ejson:ArrayGetData找不到index")
	}
	m := v[key]
	return NewEjson(m)

}

func (e *Ejson) ArrayGetMap(key int) (*Ejson, error) {

	v, ok := e.getj().([]interface{})
	if !ok {
		return nil, errors.New("ejson:arraygetmap主体对象必须是array")
	}
	if key > (len(v) - 1) {
		return nil, errors.New("ejson:arraygetmap找不到index")
	}
	m, ok := (v[key]).(map[string]interface{})
	if !ok {
		return nil, errors.New("ejson:arraygetmap找到的index不是map类型")
	}

	return NewEjson(m)
}

func (e *Ejson) ArrayGetArray(key int) (*Ejson, error) {

	v, ok := e.getj().([]interface{})
	if !ok {
		return nil, errors.New("ejson:ArrayGetArray主体对象必须是array")
	}
	if key > (len(v) - 1) {
		return nil, errors.New("ejson:ArrayGetArray找不到index")
	}
	m, ok := (v[key]).([]interface{})
	if !ok {
		return nil, errors.New("ejson:ArrayGetArray找到的index不是array类型")
	}
	return NewEjson(m)
}

func (e *Ejson) MapGetData(key string) (interface{}, error) {

	v, ok := e.getj().(map[string]interface{})
	if !ok {
		return nil, errors.New("ejson:MapGetData主体对象必须是map")
	}
	find, ok := v[key]
	if !ok {
		return nil, errors.New("ejson:MapGetData找不到index")
	}
	return find, nil
}

func (e *Ejson) MapGetMap(key string) (*Ejson, error) {

	v, ok := e.getj().(map[string]interface{})
	if !ok {
		return nil, errors.New("ejson:MapGetMap主体对象必须是map")
	}
	find, ok := v[key]
	if !ok {
		return nil, errors.New("ejson:MapGetMap找不到index")
	}
	arr, ok := find.(map[string]interface{})
	if !ok {
		return nil, errors.New("ejson:MapGetMap找到的index不是map类型")
	}
	return NewEjson(arr)
}

func (e *Ejson) ArrayGetFloat64(key int) (float64, error) {

	v, ok := e.getj().([]interface{})
	if !ok {
		return 0, errors.New("ejson:ArrayGetFloat64主体对象必须是数组")
	}
	if key > (len(v) - 1) {
		return 0, errors.New("ejson:ArrayGetFloat64找不到index")
	}
	f, ok := v[key].(float64)
	if !ok {
		return 0, errors.New("ejson:ArrayGetFloat64找到的值不是float64类型")
	}
	return f, nil
}

func (e *Ejson) ArrayGetBool(key int) (bool, error) {

	v, ok := e.getj().([]interface{})
	if !ok {
		return false, errors.New("ejson:ArrayGetBool 主体对象必须是map")
	}
	if key > (len(v) - 1) {
		return false, errors.New("ejson:ArrayGetBool :key not exist")
	}
	value, ok := v[key].(bool)
	if !ok {
		return false, errors.New("ejson:ArrayGetBool :value not bool")
	}
	return value, nil

}
func (e *Ejson) ArrayGetString(key int) (string, error) {

	v, ok := e.getj().([]interface{})
	if !ok {
		return "", errors.New("ejson:ArrayGetString主体对象必须是数组")
	}
	if key > (len(v) - 1) {
		return "", errors.New("ejson:ArrayGetString找不到index")
	}
	str, ok := v[key].(string)
	if !ok {
		return "", errors.New("ejson:ArrayGetString找到的值不是String类型")
	}
	return str, nil
}

func (e *Ejson) MapGetArray(key string) (*Ejson, error) {

	v, ok := e.getj().(map[string]interface{})
	if !ok {
		return nil, errors.New("ejson:MapGetArray主体对象必须是map")
	}
	find, ok := v[key]
	if !ok {
		return nil, errors.New(fmt.Sprintln("ejson:MapGetArray找不到index: ", key))
	}
	arr, ok := find.([]interface{})
	if !ok {
		return nil, errors.New(fmt.Sprintln("ejson:MapGetArray找到的index不是array类型", key))
	}
	return NewEjson(arr)

}

func (e *Ejson) MapGetValue(key string) (interface{}, error) {
	v, ok := e.getj().(map[string]interface{})
	if !ok {
		return 0, errors.New("ejson:MapGetFloat64主体对象必须是map")
	}
	find, ok := v[key]
	if !ok {
		return 0, errors.New("ejson:MapGetFloat64找不到index")
	}
	return find, nil
}

func (e *Ejson) MapGetFloat64(key string) (float64, error) {

	v, ok := e.getj().(map[string]interface{})
	if !ok {
		return 0, errors.New("ejson:MapGetFloat64主体对象必须是map:" + key)
	}
	find, ok := v[key]
	if !ok {
		return 0, errors.New("ejson:MapGetFloat64找不到index:" + key)
	}
	num, ok := find.(float64)
	if !ok {
		numInt, ok := find.(int)
		if !ok {
			return 0, errors.New("ejson:MapGetFloat64找到的index不是float64或者int:" + key)
		}
		num = float64(numInt)
	}
	return num, nil

}

func (e *Ejson) MapGetBool(key string) (bool, error) {

	v, ok := e.getj().(map[string]interface{})
	if !ok {
		return false, errors.New("ejson:MapGetBool主体对象必须是map")
	}
	find, ok := v[key]
	if !ok {
		return false, errors.New("ejson:MapGetBool:key not exist")
	}
	value, ok := find.(bool)
	if !ok {
		return false, errors.New("ejson:MapGetBool:value not bool")
	}
	return value, nil

}

func (e *Ejson) MapGetString(key string) (string, error) {

	v, ok := e.getj().(map[string]interface{})
	if !ok {
		return "", errors.New("ejson:MapGetString对象本身必须是map" + key)
	}
	find, ok := v[key]
	if !ok {
		return "", errors.New("ejson:MapGetString找不到index" + key)
	}
	str, ok := find.(string)
	if !ok {
		return "", errors.New("ejson:MapGetString找到的数据不是string" + key)
	}
	return str, nil

}

func (e *Ejson) IsMap() bool {

	_, ok := e.getj().(map[string]interface{})
	return ok
}

func (e *Ejson) IsArray() bool {

	_, ok := e.getj().([]interface{})
	return ok
}

//  其作用是复制一个map json, 解决在收集筛选结果的时候复制了map本身，底层数据仍然公用
//  当然，数据的复制是有性能代价的
func CopyMap(m map[string]interface{}) map[string]interface{} {
	tmp, err := NewEjson(m)
	if err == nil {
		result, err := tmp.CopyMap()
		if err == nil {
			return result
		}
	}
	return make(map[string]interface{})
}
