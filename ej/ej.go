package ej

//json



import (
	"github.com/rz1226/rzlib/errset"
	"encoding/json"
	"errors"
	"fmt"
)
/**
//不用每次检查错误，错误累积到最后一起检查
ej := NewEj("...")
data := ej.FetchFromMap("key").
                                .First()
                                .Index()

*/
type Ej struct {
	j       interface{} //json对象，一般是[]interface{}或者map[string]interface{}
	s       []byte      //json字符串
	err  *errset.ErrSet
}

func NewEj(j interface{}) *Ej   {
	e := Ej{}
	e.err  = errset.NewErrSet()
	e.err.Log("newej")
	if jdata, ok := j.(string); ok {
		e.s = []byte(jdata)
		err := json.Unmarshal(e.s, &e.j)
		if err != nil {
			e.err.Add(err )
		}
		return &e
	}
	if jdata, ok := j.([]byte); ok {
		e.s = jdata
		err := json.Unmarshal(e.s, &e.j)
		if err != nil {
			e.err.Add(err )
		}
		return &e
	}
	_, ok := j.([]interface{})
	_, ok2 := j.(map[string]interface{})
	if ok == false && ok2 == false {
		e.err.Add( errors.New("not slice or map, invalid json") )
		return &e
	}
	e.j = j
	s, err := json.Marshal(e.j)
	if err != nil {
		e.err.Add(err )
		return &e
	}
	e.s = s
	return &e
}


func (e *Ej) Json() (interface{}, error) {
	if e.err.HaveErr(){
		return nil, e.err.Error()
	}
	return e.j, nil
}

func (e *Ej) Bytes() ([]byte, error) {
	if e.err.HaveErr(){
		return nil, e.err.Error()
	}
	return e.s, nil
}

func (e *Ej) String() (string, error) {
	if e.err.HaveErr(){
		return "", e.err.Error()
	}
	return string(e.s), nil
}
// 取出一个[]map[string]interface{}, 因为这个格式非常常用做个封装
func (e *Ej) FetchDataLikeQueryRes() ([]map[string]interface{}, error ){
	if e.err.HaveErr(){
		return nil, e.err.Error()
	}

	data , ok := e.j.([]interface{})
	if !ok {
		return nil, errors.New("格式错误, json不是[]interface{}")
	}
	res := make([]map[string]interface{},0,10)
	for _, v :=range data {
		ele , ok := v.(map[string]interface{})
		if !ok {
			return nil , errors.New("json数据内部元素的格式不是map[string]interface{}")
		}
		res = append( res, ele )
	}
	return res, nil

}


func (e *Ej) FetchFromMap(key string ) ( *Ej ) {
	e.err.Log("fetchfrommap")
	if e.err.HaveErr(){
		return e
	}
	v, ok := e.j.(map[string]interface{})
	if !ok {
		e.err.Add(errors.New("fetchfrommap err :data is not map for key :"+key ))
		return e
	}
	find, ok := v[key]
	if !ok {
		e.err.Add(errors.New("fetchfrommap找不到key :"+key ))
		return e
	}
	return NewEj(find )
}

func (e *Ej) First() ( *Ej ) {
	return e.Index(0)
}

func (e *Ej) Index(key int) ( *Ej ) {
	e.err.Log("Index")
	if e.err.HaveErr(){
		return e
	}
	v, ok := e.j.([]interface{})
	if !ok {
		e.err.Add(errors.New("Index err :data is not array for key :"+ fmt.Sprint(key )))
		return e
	}
	if key > (len(v) - 1) {
		e.err.Add(errors.New("Index err :can not find key out of range:"+ fmt.Sprint(key )))
		return e
	}
	return NewEj(v[key])
}
