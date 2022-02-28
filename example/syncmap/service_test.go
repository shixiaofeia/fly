package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"testing"
)

func TestMap(t *testing.T) {
	m := Map{}
	m.Store("name", "fly")
	log.Println(m.Load("name"))
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			m.Store(strconv.Itoa(i), i)
		}(i)
	}
	for i := 0; i < 99; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			m.Delete(strconv.Itoa(i))
		}(i)
	}
	wg.Wait()
	m.Store("age", "25")
	m.Range(func(key, value interface{}) bool {
		fmt.Println(key, value)
		return true
	})
	m.Store("sex", 1)
	log.Println(m.Load("sex"))
}
