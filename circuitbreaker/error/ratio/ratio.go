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

const resName = "circuit-breaker-error-ratio"

type stateChangeTestListener struct{}

func (s *stateChangeTestListener) OnTransformToClosed(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	fmt.Printf("rule.steategy: %+v, From %s to Closed, time: %d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

func (s *stateChangeTestListener) OnTransformToOpen(prev circuitbreaker.State, rule circuitbreaker.Rule, snapshot interface{}) {
	fmt.Printf("rule.steategy: %+v, From %s to Open, snapshot: %.2f, time: %d\n", rule.Strategy, prev.String(), snapshot, util.CurrentTimeMillis())
}

func (s *stateChangeTestListener) OnTransformToHalfOpen(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	fmt.Printf("rule.steategy: %+v, From %s to Half-Open, time: %d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

func main() {
	err := sentinel.InitWithConfigFile("./sentinel.yml")
	if err != nil {
		log.Fatal(err)
	}
	ch := make(chan struct{})

	circuitbreaker.RegisterStateChangeListeners(&stateChangeTestListener{})

	_, err = circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		{
			Resource:                     resName,
			Strategy:                     circuitbreaker.ErrorRatio,
			RetryTimeoutMs:               3000,  // 熔断后3s内，请求快速失败
			MinRequestAmount:             10,    // 静默数量，若当前统计周期内的请求数小于此值，即使达到熔断条件规则也不会触发
			StatIntervalMs:               60000, // 统计的时间窗口大小60s
			StatSlidingWindowBucketCount: 10,    // StatIntervalMs%StatSlidingWindowBucketCount == 0
			Threshold:                    0.3,   // 窗口期内触发熔断error占比
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
		e, b := sentinel.Entry(resName)
		if b != nil {
			// g2 blocked
			fmt.Println("blocked")
		} else {
			if rand.Uint64()%20 > 9 {
				fmt.Println("error")
				e.SetError(errors.New("error"))
			} else {
				fmt.Println("success")
			}
			// g2 passed
			e.Exit()
		}
		time.Sleep(time.Millisecond * 500)
	}
}
