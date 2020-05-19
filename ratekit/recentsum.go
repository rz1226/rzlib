package ratekit

/*
计算最近收集的数字的和
*/

import (
	"sync"
)

type RecentSum struct {
	ints []int // 这个结构的目的就是计算这里的和
	len  int   // 只保留len个数字，多余的扔出去
	mu   *sync.RWMutex
}

func NewRecentSum(length int) *RecentSum {
	if length <= 0 {
		length = 5
	}
	r := &RecentSum{}
	r.len = length
	r.mu = &sync.RWMutex{}
	r.ints = make([]int, length)
	for i := 0; i < length; i++ {
		r.ints[i] = 0
	}
	return r
}

func (r *RecentSum) Put(n int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ints = append(r.ints[1:r.len], n)

}
func (r *RecentSum) Sum() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	sum := 0
	for _, e := range r.ints {
		sum += e
	}
	return sum
}
