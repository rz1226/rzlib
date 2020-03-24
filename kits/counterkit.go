package kits

type CounterKit struct {
	data   *Counter
	name   string
	readme string //名称的注释
}

func NewCounterKit(name, readme string) *CounterKit {
	c := &CounterKit{}
	c.name = name
	c.data = NewCounter()
	c.readme = readme
	return c
}

func (c *CounterKit) Show() string {
	str := ""
	str += "\n----------------------\n计数器名称:" + c.name + " : \n计数器信息:" + c.readme + "\n"
	str += c.Str()
	str += "\n"

	return str
}

func (c *CounterKit) Inc() {

	c.data.Add(1)
}
func (c *CounterKit) IncBy(num int64) {

	c.data.Add(num)
}
func (c *CounterKit) Get(name string) int64 {
	return c.data.Get()
}

func (c *CounterKit) Str() string {
	return c.data.Str()
}
