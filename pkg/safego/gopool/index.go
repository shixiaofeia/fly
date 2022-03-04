package gopool

import "sync"

type Pool struct {
	wg    sync.WaitGroup
	queue chan struct{}
}

// NewGoPool 实例化一个go程池
func NewGoPool(i int) *Pool {
	if i < 1 {
		i = 1
	}
	return &Pool{queue: make(chan struct{}, i)}
}

// Add 添加
func (p *Pool) Add() {
	p.queue <- struct{}{}
	p.wg.Add(1)
}

// Done 释放
func (p *Pool) Done() {
	p.wg.Done()
	<-p.queue
}

// Wait 等待
func (p *Pool) Wait() {
	p.wg.Wait()
}
