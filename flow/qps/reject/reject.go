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

func main() {
	// config https://sentinelguard.io/zh-cn/docs/golang/general-configuration.html
	err := sentinel.InitWithConfigFile("./sentinel.yml")
	if err != nil {
		log.Fatal(err)
	}

	// 配置一条限流规则 https://sentinelguard.io/zh-cn/docs/golang/flow-control.html
	// 可以通过动态文件、etcd、consul 等配置中心来动态地配置规则。
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "some-test",
			Threshold:              10,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			StatIntervalInMs:       1000,

			// ControlBehavior:        flow.Throttling,
			// MaxQueueingTimeMs:      500, // 排队最大等待时间，平滑流量波动，更好应对脉冲流量

			// WarmUpPeriodSec:        10,  // 预热时间
			// WarmUpColdFactor:       3,   // 预热因子，默认是 3
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	ch := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func() {
			for {
				// 埋点逻辑，埋点资源名为 some-test
				e, b := sentinel.Entry("some-test")
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
