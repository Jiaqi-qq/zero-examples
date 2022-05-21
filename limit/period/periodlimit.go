package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

const seconds = 5

var (
	rdx     = flag.String("redis", "localhost:6379", "the redis, default localhost:6379")
	rdxPass = flag.String("redisPass", "", "the redis password")
	rdxKey  = flag.String("redisKey", "rate", "the redis key, default rate")
	threads = flag.Int("threads", runtime.NumCPU(), "the concurrent threads, default to cores")
)

func main() {
	flag.Parse()

	store := redis.New(*rdx, redis.WithPass(*rdxPass))
	fmt.Println(store.Ping())
	lmt := limit.NewPeriodLimit(seconds, 3, store, *rdxKey) // 每5秒允许3次
	timer := time.NewTimer(time.Second * seconds * 10)      // 延长程序运行时间
	quit := make(chan struct{})
	defer timer.Stop()
	go func() {
		<-timer.C
		close(quit)
	}()

	var allowed, denied int32
	var wait sync.WaitGroup
	for i := 0; i < *threads; i++ {
		i := i
		wait.Add(1)
		go func() {
			for {
				time.Sleep(time.Second)
				select {
				case <-quit:
					wait.Done()
					return
				default:
					key := strconv.FormatInt(int64(i), 10)
					v, err := lmt.Take(key)
					if err == nil && v == limit.Allowed {
						fmt.Printf("AllowedQuota key: %v\n", key)
						atomic.AddInt32(&allowed, 1)
					} else if err != nil {
						log.Fatal(err)
					} else if v == limit.OverQuota {
						fmt.Printf("OverQuota key: %v\n", key)
						atomic.AddInt32(&denied, 1)
					} else if v == limit.HitQuota {
						fmt.Printf("HitQuota key: %v\n", key)
					}
				}
			}
		}()
	}

	wait.Wait()
	fmt.Printf("allowed: %d, denied: %d, qps: %d\n", allowed, denied, (allowed+denied)/seconds)
}
