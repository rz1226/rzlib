package cache

//内存lru cache

import (
	"errors"
	"github.com/karlseguin/ccache"
	"time"
)
type CCache struct{
	c *ccache.Cache
}
func NewCCache(maxSize int ) *CCache{
	cache := new(CCache)
	count := uint32(maxSize/10 + 1 )
	cache.c = ccache.New(ccache.Configure().MaxSize( int64(maxSize) ).ItemsToPrune(count))
	return cache
}


func (cd Data)ToCCache(source *CCache , duration time.Duration ) error {
	if len(cd.key) == 0 {
		return errors.New("no key for ccache")
	}
	source.c.Set(cd.key, cd.data, duration)
	return nil
}

func (k DataKey) FetchFromCCache(source *CCache)( interface{}, error) {
	key := string(k )
	item := source.c.Get(key )
	if item != nil {
		if item.TTL().Seconds() > 0 {
			return item.Value(), nil
		}
	}
	return nil, errors.New("no data")
}
