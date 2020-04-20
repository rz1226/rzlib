package gobutil

import (
	"encoding/gob"
	"os"
)

func Save(path string, object interface{}) error {
	file, err := os.Create(path)
	if err == nil {
		defer file.Close()
		encoder := gob.NewEncoder(file)
		return  encoder.Encode(object)
	}
	return err
}

func Load(path string, object interface{}) error {
	file, err := os.Open(path)
	if err == nil {
		defer file.Close()
		decoder := gob.NewDecoder(file)
		return  decoder.Decode(object)
	}
	return err
}

