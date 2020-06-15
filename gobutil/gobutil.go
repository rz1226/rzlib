package gobutil

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"os"
)

func Save(path string, object interface{}) error {
	file, err := os.Create(path)
	if err == nil {
		defer file.Close()
		encoder := gob.NewEncoder(file)
		return encoder.Encode(object)
	}
	return err
}

//第二个参数是指针
func Load(path string, object interface{}) error {
	file, err := os.Open(path)
	if err == nil {
		defer file.Close()
		decoder := gob.NewDecoder(file)
		return decoder.Decode(object)
	}
	return err
}

/*****************************/

func ToBytes(object interface{}) ([]byte, error) {
	w := bytes.NewBuffer(nil)

	encoder := gob.NewEncoder(w)
	err := encoder.Encode(object)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(w)

}

// 参数是指针
func ToStruct(data []byte, object interface{}) error {
	reader := bytes.NewReader(data)
	decoder := gob.NewDecoder(reader)
	err := decoder.Decode(object)
	if err != nil {
		return err
	}
	return nil
}
