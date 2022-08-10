package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/util"
)

type stateChangeTestListener struct{}

func (s *stateChangeTestListener) OnTransformToClosed(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	fmt.Printf("rule.steategy: %+v, From %s to Closed, time: %d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

func (s *stateChangeTestListener) OnTransformToOpen(prev circuitbreaker.State, rule circuitbreaker.Rule, snapshot interface{}) {
	fmt.Printf("rule.steategy: %+v, From %s to Open, snapshot: %d, time: %d\n", rule.Strategy, prev.String(), snapshot, util.CurrentTimeMillis())
}

func (s *stateChangeTestListener) OnTransformToHalfOpen(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	fmt.Printf("rule.steategy: %+v, From %s to Half-Open, time: %d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

func main() {
	// conf := config.NewDefaultConfig()
	// conf.Sentinel.Log.Logger = logging.NewConsoleLogger()
	// err := sentinel.InitWithConfig(conf)
	err := sentinel.InitWithConfigFile("./sentinel.yml")
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	ch := make(chan struct{})

	circuitbreaker.RegisterStateChangeListeners(&stateChangeTestListener{})

	_, err = circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		{
			Resource:                     "abc",
			Strategy:                     circuitbreaker.ErrorCount,
			RetryTimeoutMs:               3000,  // 熔断后3s内，请求快速失败
			MinRequestAmount:             10,    // 静默数量，若当前统计周期内的请求数小于此值，即使达到熔断条件规则也不会触发
			StatIntervalMs:               60000, // 统计的时间窗口大小60s
			StatSlidingWindowBucketCount: 10,    // StatIntervalMs%StatSlidingWindowBucketCount == 0
			Threshold:                    5,     // 窗口期内触发熔断error次数
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	go handle()
	go handle()

	<-ch
}

func handle() {
	for {
		e, b := sentinel.Entry("abc")
		if b != nil {
			fmt.Println("blocked")
		} else {
			if rand.Uint64()%20 > 9 {
				fmt.Println("error")
				e.SetError(errors.New("error"))
			} else {
				fmt.Println("success")
			}
			e.Exit()
		}
		time.Sleep(time.Millisecond * 500)
	}
}
