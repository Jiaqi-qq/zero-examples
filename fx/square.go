package main

import (
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/fx"
)

func main() {
	result, err := fx.From(func(source chan<- interface{}) {
		for i := 0; i < 5; i++ {
			source <- i
			source <- i
		}
	}).Map(func(item interface{}) interface{} {
		i := item.(int)
		return i * i
	}).Filter(func(item interface{}) bool {
		i := item.(int)
		return i%2 == 0
	}).Distinct(func(item interface{}) interface{} {
		fmt.Println("distinct: ", item)
		return item
	}).Reduce(func(pipe <-chan interface{}) (interface{}, error) {
		time.Sleep(time.Second * 20)
		var result int
		for item := range pipe {
			fmt.Println("Reduce: ", item)
			i := item.(int)
			result += i
		}
		return result, nil
	})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}
}
