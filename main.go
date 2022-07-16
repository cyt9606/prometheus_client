package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/push"
	"time"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/node_exporter/collector"
)








type collectorAdapter struct {
	collector.Collector
}

// Describe implements the prometheus.Collector interface.
func (a collectorAdapter) Describe(ch chan<- *prometheus.Desc) {
	// We have to send *some* metric in Describe, but we don't know which ones
	// we're going to get, so just send a dummy metric.
	ch <- prometheus.NewDesc("dummy_metric", "Dummy metric.", nil, nil)
}

// Collect implements the prometheus.Collector interface.
func (a collectorAdapter) Collect(ch chan<- prometheus.Metric) {
	if err := a.Update(ch); err != nil {
		panic("failed to update collector")
	}
}


func PushToGateway(pushad string,coldir string,job string,coltime int) {

	mtime := 1.0
	c := &textFileCollector{
		path:   coldir,
		mtime:  &mtime,
	}


	registry := prometheus.NewRegistry()
	registry.MustRegister(collectorAdapter{c})
	fmt.Println("collector started\n")
	pushgate:=push.New(pushad,job)
	pushgate.Gatherer(registry)
	fmt.Println("start push\n")
	for {
		if err := pushgate.Push(); err != nil {
			print(err)
		}

		time.Sleep(time.Duration(coltime)*time.Second)
	}

}
func main(){
	PushToGateway()
}
