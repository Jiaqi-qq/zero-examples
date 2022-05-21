package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/stringx"
	"github.com/zeromicro/go-zero/core/syncx"
)

func main() {
	const round = 5
	var wg sync.WaitGroup
	barrier := syncx.NewSingleFlight()

	wg.Add(round)
	for i := 0; i < round; i++ {
		go func() {
			defer wg.Done()
			val, b, err := barrier.DoEx("once", func() (interface{}, error) {
				time.Sleep(time.Second * 3)
				return stringx.RandId(), nil
			})
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(val, b)
			}
		}()
	}

	wg.Wait()
}
