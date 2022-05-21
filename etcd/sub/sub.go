package main

import (
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
)

func main() {
	sub, err := discov.NewSubscriber([]string{":2379"}, "028F2C35852D", discov.Exclusive())
	//sub, err := discov.NewSubscriber([]string{"etcd.discovery:2379"}, "028F2C35852D", discov.Exclusive())
	logx.Must(err)

	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()

	sub.AddListener(func() {
		fmt.Println("listener in ")
	})

	for {
		select {
		case <-ticker.C:
			fmt.Println("values:", sub.Values())
		}
	}
}
