package main

import (
	"log"
	"math/rand"
	"time"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/system"
	"github.com/alibaba/sentinel-golang/core/system_metric"
)

const resName = "system-load"

func main() {
	err := sentinel.InitWithConfigFile("./sentinel.yml")
	if err != nil {
		log.Fatal(err)
	}

	_, err = system.LoadRules([]*system.Rule{
		{
			MetricType:   system.Load,
			TriggerCount: 8.0,
			Strategy:     system.BBR,
		},
	})
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
		return
	}

	// mock memory usage is 1000 bytes, so QPS threshold should be 1000
	system_metric.SetSystemMemoryUsage(999)
	ch := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func() {
			for {
				e, b := sentinel.Entry(resName, sentinel.WithTrafficType(base.Inbound))
				if b != nil {
					// Blocked. We could get the block reason from the BlockError.
					time.Sleep(time.Duration(rand.Uint64()%2) * time.Millisecond)
				} else {
					// Passed, wrap the logic here.
					time.Sleep(time.Duration(rand.Uint64()%2) * time.Millisecond)
					// Be sure the entry is exited finally.
					e.Exit()
				}
			}
		}()
	}
	<-ch
}
