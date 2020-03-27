package ejson

import (
	"errors"
	"fmt"
	"strconv"
	//"reflect"
)

//比较是否相等， 字符串和数字
func compareString(a, b interface{}) bool {
	v1, err := ConvertToString(a)
	if err != nil {
		return false
	}

	v2, err := ConvertToString(b)
	if err != nil {
		return false
	}
	fmt.Println(v1, v2)
	return v1 == v2
}

func ConvertToStringNoError(input interface{}) string {
	s, err := ConvertToString(input)
	if err != nil {
		s = ""
	}
	return s
}

func ConvertToString(input interface{}) (string, error) {

	switch val := input.(type) {
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64), nil
	case int:
		return strconv.Itoa(val), nil
	case string:
		return val, nil
	case byte:
		return string(val), nil
	case bool:
		if val == true {
			return "true", nil
		} else {
			return "false", nil
		}
	default:
		return "", errors.New("convert to string error")
	}

}
