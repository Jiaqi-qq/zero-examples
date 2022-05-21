package main

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/zeromicro/go-zero/core/fx"
)

func TestDistinct(t *testing.T) {
	fx.Just(1, 3, 5, 7, 1, 3, 5, 5, 5, 5).Distinct(func(item interface{}) interface{} {
		// 通过key进行去重，相同key只保留一个
		return item
	}).ForEach(func(item interface{}) {
		fmt.Println(item)
	})
}

func TestFilter(t *testing.T) {
	fx.Just(1, 2, 3, 4, 5, 6, 7, 8, 9, 10).Filter(func(item interface{}) bool {
		if item.(int)%3 == 0 {
			return true
		}
		return false
	}).ForEach(func(item interface{}) {
		fmt.Println(item)
	})
}

func TestGroup(t *testing.T) {
	fx.Just(1, 3, 5, 7, 1, 3, 5, 5, 5, 5).Group(func(item interface{}) interface{} {
		return item
	}).ForEach(func(item interface{}) {
		fmt.Println(item)
	})
}

func TestMap(t *testing.T) {
	fx.Just(1, 3, 5).Map(func(item interface{}) interface{} {
		return strconv.Itoa(item.(int)) + " - abc"
	}).ForEach(func(item interface{}) {
		fmt.Println(item)
	})
}

func TestMerge(t *testing.T) {
	fx.Just("abc", 3, "dtf", 5).Merge().ForEach(func(item interface{}) {
		fmt.Println(item)
	})
}

func TestSkip(t *testing.T) {
	fx.Just("abc", 3, "dtf", 5).Split(2).ForEach(func(item interface{}) {
		fmt.Println(item)
	})
}

func TestSplit(t *testing.T) {
	fx.Just(1, 2, 3, 4, 5, 6, 7, 8, 9, 10).Split(4).ForEach(func(item interface{}) {
		vals := item.([]interface{})
		fmt.Println(len(vals), vals)
	})
}

func TestForAll(t *testing.T) {
	fx.Just(1, 2, 3, 4).ForAll(func(pipe <-chan interface{}) {
		for x := range pipe {
			fmt.Println(x)
		}
	})
}

func TestParallel(t *testing.T) {
	fx.Just(1, 2, 3, 4, 5, 6, 7, 8, 9, 10).Group(func(item interface{}) interface{} {
		return item.(int) % 3
	}).Parallel(func(item interface{}) {
		fmt.Println(item)
		time.Sleep(time.Second * 2)
	})
}

func TestReduce(t *testing.T) {
	reduce, err := fx.Just(1, 2, 3, 4).Reduce(func(pipe <-chan interface{}) (interface{}, error) {
		for x := range pipe {
			time.Sleep(time.Second)
			fmt.Println(x)
		}
		return "你以为这是什么", nil
	})
	fmt.Println(reduce, err)
}

func BenchmarkFx(b *testing.B) {
	type Mixed struct {
		Name   string
		Age    int
		Gender int
	}
	for i := 0; i < b.N; i++ {
		var mx Mixed
		fx.Parallel(func() {
			mx.Name = "hello"
		}, func() {
			mx.Age = 20
		}, func() {
			mx.Gender = 1
		})
	}
}
