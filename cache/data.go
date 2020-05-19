package cache

/*
a := "some string"
NewData( a ).SetKey(key).ToCCache( Source, seconds )


data :=  NewKey( key ).FetchFromCCache( Source )

*/

type Data struct {
	data interface{} //  数据本身
	key  string      //  可用于  kv  key,    hash key
}

func NewData(data interface{}) *Data {
	d := new(Data)
	d.data = data
	return d
}

func (d *Data) SetKey(key string) *Data {
	d.key = key
	return d
}

type DataKey string

func NewKey(key string) DataKey {
	return DataKey(key)
}
