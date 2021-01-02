package engine

import (
	"fmt"
	"sync"
	"time"

	"github.com/mauriliommachado/go-metric/storage"
)

// WorkerAPI asdas
type WorkerAPI struct {
	input           chan storage.Event
	output          chan storage.Event
	results         chan int
	quit            chan int
	Storage         *storage.API
	deltaSlice      *DeltaSlice
	numberOfWorkers int
}

//New todo
func New(numberOfWorkers int) *WorkerAPI {
	var wapi WorkerAPI
	wapi.numberOfWorkers = numberOfWorkers
	wapi.input = make(chan storage.Event, 10)
	wapi.output = make(chan storage.Event, 10)
	wapi.results = make(chan int, 10)
	wapi.quit = make(chan int)
	wapi.Storage = storage.New()
	wapi.newDeltaSlice(60000)
	for i := 1; i <= numberOfWorkers; i++ {
		go inputWorker(i, &wapi)
	}
	for i := 1; i <= numberOfWorkers; i++ {
		go outputWorker(i, &wapi)
	}
	fmt.Println("Metric store initialized")
	return &wapi
}

func inputWorker(id int, wapi *WorkerAPI) {
	fmt.Println("Input Worker", id, " waiting for work")
	for {
		select {
		case event := <-wapi.input:
			ok := wapi.deltaSlice.Push(event)
			if ok {
				wapi.do(event, true)
			}
		case <-wapi.quit:
			fmt.Println("quit")
			return
		}
	}
}

func outputWorker(id int, wapi *WorkerAPI) {
	fmt.Println("Output Worker", id, " waiting for work")
	for {
		select {
		case event := <-wapi.output:
			wapi.do(event, false)
		case <-wapi.quit:
			fmt.Println("quit")
			return
		}
	}
}

func (wapi *WorkerAPI) do(event storage.Event, inputEvent bool) {
	var waitgroup sync.WaitGroup
	for _, metric := range wapi.Storage.GetMetricDefinitions() {
		switch {
		case metric.MetricType == "count":
			waitgroup.Add(1)
			go wapi.doCount(metric, event, inputEvent, &waitgroup)
		case metric.MetricType == "sum":
			waitgroup.Add(1)
			go wapi.doSum(metric, event, inputEvent, &waitgroup)
		case metric.MetricType == "max":
			waitgroup.Add(1)
			go wapi.doMax(metric, event, inputEvent, &waitgroup)
		case metric.MetricType == "min":
			waitgroup.Add(1)
			go wapi.doMin(metric, event, inputEvent, &waitgroup)
		default:
			fmt.Println("Metric type unkown")
		}
	}
	waitgroup.Wait()
}

func updatePartialSlices() {

}

//SubmitWork todo
func (wapi *WorkerAPI) SubmitWork(data map[string]interface{}) {
	var event storage.Event
	event.Data = data
	if val, ok := data["timestamp"]; ok {
		event.Timestamp = int64(val.(float64))
	} else {
		event.Timestamp = nowAsUnixMilli()
	}
	wapi.input <- event
}

func nowAsUnixMilli() int64 {
	return time.Now().UnixNano() / 1e6
}
