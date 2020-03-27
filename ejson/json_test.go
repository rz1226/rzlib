package ejson

import (
	"fmt"
	"testing"
)

type Data struct {
	Abc string
	Cd  int
	X   string
}

//测试map到struct的映射
func Test_mapstruct(t *testing.T) {
	source := make(map[string]interface{})
	source["abc"] = "abcxxx"
	source["cd"] = 12
	source["t"] = 34
	data := &Data{}
	err := MapToStruct(data, source)
	fmt.Println(err)
	fmt.Println(data)
}
