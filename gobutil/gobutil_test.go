package gobutil

import (
	"fmt"
	"os"
	"runtime"
	"testing"
)

const file = "./test.gob"

type User struct {
	Name, Pass string
}

func Test_save_load(t *testing.T) {
	var datato = &User{"Donald", "DuckPass"}
	var datafrom = new(User)

	err := Save(file, datato)
	Check(err)
	err = Load(file, datafrom)
	Check(err)
	fmt.Println(datafrom)

}

func Test_save_load_map(t *testing.T) {
	var datato = map[string]string{"a": "b"}
	var datafrom = make(map[string]string)

	err := Save(file, datato)
	Check(err)
	err = Load(file, &datafrom)
	Check(err)
	fmt.Println(datafrom)

}

func Check(e error) {
	if e != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Println(line, "\t", file, "\n", e)
		os.Exit(1)
	}
}
