package rztool

import (
	"bytes"
	"io/ioutil"
)

//  处理文件

type FilePath string

//  读取文件内容
func (fp FilePath) GetContent() (string, error) {
	res, err := ioutil.ReadFile(string(fp))
	if err != nil {
		return "", err
	}
	//  去掉bom
	res = bytes.TrimPrefix(res, []byte("\xef\xbb\xbf"))
	return string(res), nil
}
