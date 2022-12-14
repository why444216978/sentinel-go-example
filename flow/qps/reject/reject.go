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

const resName = "flow-qps-reject"

func main() {
	// config https://sentinelguard.io/zh-cn/docs/golang/general-configuration.html
	// conf := config.NewDefaultConfig()
	// conf.Sentinel.Log.Logger = logging.NewConsoleLogger()
	// err := sentinel.InitWithConfig(conf)
	err := sentinel.InitWithConfigFile("./sentinel.yml")
	if err != nil {
		log.Fatal(err)
	}

	// 配置一条限流规则 https://sentinelguard.io/zh-cn/docs/golang/flow-control.html
	// 可以通过动态文件、etcd、consul 等配置中心来动态地配置规则。
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               resName,
			Threshold:              10,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			StatIntervalInMs:       1000, // 1000ms is QPS
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
