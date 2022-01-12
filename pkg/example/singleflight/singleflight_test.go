package singleflight

import (
	"fmt"
	"go.uber.org/atomic"
	"sync"
	"testing"
	"time"
)

type City struct {
	Name string
}

var (
	cacheCity interface{}
	dbNum     atomic.Uint32
)

func TestNewSingleFlight(t *testing.T) {
	now := time.Now()
	var (
		wg = sync.WaitGroup{}
	)
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			v, err := getData("cache:city")
			if err != nil {
				t.Error(err)
			}
			t.Log(v)
		}()
	}
	wg.Wait()
	t.Log("from db ", dbNum)
	t.Log(time.Since(now))
}

// getData 获取数据
func getData(key string) (val interface{}, err error) {
	if val, err = getCityFromCache(key); err != nil {
		// 使用单程
		val, err, _ = NewSingleFlight().Do(key, func() (interface{}, error) {
			return getDbAndSetCache(key)
		})
		//正常缓存
		//return getDbAndSetCache(key)
		return
	}
	return
}

// getDbAndSetCache 从db获取信息并且同步到缓存
func getDbAndSetCache(key string) (val interface{}, err error) {
	if val, err = getCityFromDb(key); err == nil {
		// set cache
		cacheCity = val
	}
	return val, err
}

// getCityFromDb 从db获取信息
func getCityFromDb(key string) (interface{}, error) {
	dbNum.Add(1)
	cityList := make([]City, 0)
	for _, v := range []string{"北京", "上海"} {
		cityList = append(cityList, City{Name: v})
	}
	// 模拟数据库查询时间
	time.Sleep(10 * time.Millisecond)
	return cityList, nil
}

// getCityFromCache 从缓存获取信息
func getCityFromCache(key string) (interface{}, error) {
	if cacheCity == nil {
		return "", fmt.Errorf("not found")
	}
	return cacheCity, nil
}
