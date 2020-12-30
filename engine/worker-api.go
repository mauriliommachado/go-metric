package engine

import (
	"fmt"
	"math/bits"
	"strconv"

	"github.com/mauriliommachado/go-metric/storage"
)

// WorkerAPI asdas
type WorkerAPI struct {
	input           chan map[string]interface{}
	results         chan int
	quit            chan int
	Storage         *storage.StorageAPI
	numberOfWorkers int
}

//New todo
func New(numberOfWorkers int) *WorkerAPI {
	var wapi WorkerAPI
	wapi.numberOfWorkers = numberOfWorkers
	wapi.input = make(chan map[string]interface{}, 100)
	wapi.results = make(chan int, 10)
	wapi.quit = make(chan int)
	wapi.Storage = storage.New()
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
			fmt.Println(event, " data received on worker ", id)
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
		case <-wapi.quit:
			fmt.Println("quit")
			return
		}
	}
}

//SubmitWork todo
func (wapi *WorkerAPI) SubmitWork(data map[string]interface{}) {
	wapi.input <- data
}

func (wapi *WorkerAPI) doMax(metric storage.MetricDefinition, event map[string]interface{}) {
	key := metric.GetDefinitionKey(event)
	metricValue, ok := wapi.Storage.Get(key)
	if !ok {
		metricValue = storage.MetricItem{Key: key, Value: (1 << bits.UintSize) / -2}
	}
	eventValue, _ := strconv.Atoi(fmt.Sprintf("%v", event[metric.Aggregation]))
	if eventValue > metricValue.Value {
		metricValue.Value = eventValue
		//saving calculated metric value
		wapi.Storage.Put(key, metricValue)
	}
}

func (wapi *WorkerAPI) doMin(metric storage.MetricDefinition, event map[string]interface{}) {
	key := metric.GetDefinitionKey(event)
	metricValue, ok := wapi.Storage.Get(key)
	if !ok {
		metricValue = storage.MetricItem{Key: key, Value: (1<<bits.UintSize)/2 - 1}
	}
	eventValue, _ := strconv.Atoi(fmt.Sprintf("%v", event[metric.Aggregation]))
	if eventValue < metricValue.Value {
		metricValue.Value = eventValue
		//saving calculated metric value
		wapi.Storage.Put(key, metricValue)
	}
}

func (wapi *WorkerAPI) doCount(metric storage.MetricDefinition, event map[string]interface{}) {
	key := metric.GetDefinitionKey(event)
	metricValue, ok := wapi.Storage.Get(key)
	if !ok {
		metricValue = storage.MetricItem{Key: key, Value: 0}
	}
	metricValue.Value++
	//saving calculated metric value
	wapi.Storage.Put(key, metricValue)
}

func (wapi *WorkerAPI) doSum(metric storage.MetricDefinition, event map[string]interface{}) {
	key := metric.GetDefinitionKey(event)
	metricValue, ok := wapi.Storage.Get(key)
	if !ok {
		metricValue = storage.MetricItem{Key: key, Value: 0}
	}
	aggregationValue, _ := strconv.Atoi(fmt.Sprintf("%v", event[metric.Aggregation]))
	metricValue.Value = metricValue.Value + aggregationValue
	//saving calculated metric value
	wapi.Storage.Put(key, metricValue)
}
