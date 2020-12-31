package engine

import (
	"fmt"
	"math/bits"
	"strconv"
	"time"

	"github.com/mauriliommachado/go-metric/storage"
)

// WorkerAPI asdas
type WorkerAPI struct {
	input           chan storage.Event
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
	wapi.input = make(chan storage.Event, 100)
	wapi.results = make(chan int, 10)
	wapi.quit = make(chan int)
	wapi.Storage = storage.New()
	wapi.newDeltaSlice(2)
	for i := 1; i <= numberOfWorkers; i++ {
		go worker(i, &wapi)
	}
	fmt.Println("Metric store initialized")
	return &wapi
}

func worker(id int, wapi *WorkerAPI) {
	fmt.Println("Worker", id, " waiting for work")
	for {
		select {
		case event := <-wapi.input:
			//fmt.Println(event, " data received on worker ", id)
			ok := wapi.deltaSlice.Push(event)
			if ok {
				wapi.do(event)
			}
			//fmt.Println("----------------------")
		case <-wapi.quit:
			fmt.Println("quit")
			return
		}
	}
}

func (wapi *WorkerAPI) do(event storage.Event) {
	for _, metric := range wapi.Storage.GetMetricDefinitions() {
		switch {
		case metric.MetricType == "count":
			wapi.doCount(metric, event)
		case metric.MetricType == "sum":
			wapi.doSum(metric, event)
		case metric.MetricType == "max":
			wapi.doMax(metric, event)
		case metric.MetricType == "min":
			wapi.doMin(metric, event)
		default:
			fmt.Println("Metric type unkown")
		}
	}
}

func (wapi *WorkerAPI) undo(event storage.Event) {
	for _, metric := range wapi.Storage.GetMetricDefinitions() {
		switch {
		case metric.MetricType == "count":
			wapi.doUnCount(metric, event)
		default:
			//fmt.Println("Metric type unkown")
		}
	}
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

func (wapi *WorkerAPI) doMax(metric storage.MetricDefinition, event storage.Event) {
	key := metric.GetDefinitionKey(event)
	metricValue, ok := wapi.Storage.Get(key)
	if !ok {
		metricValue = storage.MetricItem{Key: key, Value: (1 << bits.UintSize) / -2}
	}
	eventValue, _ := strconv.Atoi(fmt.Sprintf("%v", event.Data[metric.Aggregation]))
	if eventValue > metricValue.Value {
		metricValue.Value = eventValue
		//saving calculated metric value
		wapi.Storage.Put(key, metricValue)
	}
}

func (wapi *WorkerAPI) doMin(metric storage.MetricDefinition, event storage.Event) {
	key := metric.GetDefinitionKey(event)
	metricValue, ok := wapi.Storage.Get(key)
	if !ok {
		metricValue = storage.MetricItem{Key: key, Value: (1<<bits.UintSize)/2 - 1}
	}
	eventValue, _ := strconv.Atoi(fmt.Sprintf("%v", event.Data[metric.Aggregation]))
	if eventValue < metricValue.Value {
		metricValue.Value = eventValue
		//saving calculated metric value
		wapi.Storage.Put(key, metricValue)
	}
}

func (wapi *WorkerAPI) doCount(metric storage.MetricDefinition, event storage.Event) {
	key := metric.GetDefinitionKey(event)
	metricValue, ok := wapi.Storage.Get(key)
	if !ok {
		metricValue = storage.MetricItem{Key: key, Value: 0}
	}
	metricValue.Value++
	//fmt.Println(key, " metric value after", metricValue)
	//saving calculated metric value
	wapi.Storage.Put(key, metricValue)
}

func (wapi *WorkerAPI) doUnCount(metric storage.MetricDefinition, event storage.Event) {
	key := metric.GetDefinitionKey(event)
	fmt.Println(key)
	metricValue, _ := wapi.Storage.Get(key)
	metricValue.Value--
	//fmt.Println(key, " metric value after", metricValue)
	//saving calculated metric value
	wapi.Storage.Put(key, metricValue)
}

func (wapi *WorkerAPI) doSum(metric storage.MetricDefinition, event storage.Event) {
	key := metric.GetDefinitionKey(event)
	metricValue, ok := wapi.Storage.Get(key)
	if !ok {
		metricValue = storage.MetricItem{Key: key, Value: 0}
	}
	aggregationValue, _ := strconv.Atoi(fmt.Sprintf("%v", event.Data[metric.Aggregation]))
	metricValue.Value = metricValue.Value + aggregationValue
	//saving calculated metric value
	wapi.Storage.Put(key, metricValue)
}
