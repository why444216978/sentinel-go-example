package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/util"
)

const resName = "flow-qps-warm-up"

// warm up doc https://github.com/alibaba/Sentinel/wiki/%E9%99%90%E6%B5%81---%E5%86%B7%E5%90%AF%E5%8A%A8
func main() {
	err := sentinel.InitWithConfigFile("./sentinel.yml")
	if err != nil {
		log.Fatal(err)
	}

	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               resName,
			Threshold:              10,
			TokenCalculateStrategy: flow.WarmUp,
			ControlBehavior:        flow.Reject,
			StatIntervalInMs:       1000, // 1000ms is QPS
			WarmUpPeriodSec:        10,   // 预热时间
			WarmUpColdFactor:       3,    // 预热因子，默认是 3
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func() {
			for {
				e, b := sentinel.Entry(resName)
				if b != nil {
					// 请求被流控，可以从 BlockError 中获取限流详情
					// block 后不需要进行 Exit()
					time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
				} else {
					// 请求可以通过，在此处编写您的业务逻辑
					// 务必保证业务逻辑结束后 Exit
					fmt.Println(util.CurrentTimeMillis(), "Passed")
					time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
					e.Exit()
				}

			}
		}()
	}
	<-ch
}
