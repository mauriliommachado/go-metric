package engine

import (
	"fmt"
	"math/bits"
	"strconv"
	"sync"

	"github.com/mauriliommachado/go-metric/storage"
)

func (wapi *WorkerAPI) doMax(metric storage.MetricDefinition, event storage.Event, inputEvent bool, wg *sync.WaitGroup) {
	defer wg.Done()
	if !inputEvent {
		return
	}
	key := metric.GetDefinitionKey(event)
	metricValue, ok := wapi.Storage.Get(key)
	eventValue, _ := strconv.Atoi(fmt.Sprintf("%v", event.Data[metric.Aggregation]))
	if inputEvent {
		if !ok {
			metricValue = storage.MetricItem{Key: key, Value: (1 << bits.UintSize) / -2}
		}
		if eventValue > metricValue.Value {
			metricValue.Value = eventValue
			wapi.Storage.Put(key, metricValue)
		}
	} else if eventValue == metricValue.Value { // check if the event leaving the window is the current max
		partialValue := (1 << bits.UintSize) / -2
		for _, item := range wapi.deltaSlice.Items {
			itemValue, _ := strconv.Atoi(fmt.Sprintf("%v", item.Data[metric.Aggregation]))
			if itemValue > partialValue && itemValue != metricValue.Value {
				partialValue = itemValue
			}
		}
		metricValue.Value = partialValue
		wapi.Storage.Put(key, metricValue)
	}

}

func (wapi *WorkerAPI) doMin(metric storage.MetricDefinition, event storage.Event, inputEvent bool, wg *sync.WaitGroup) {
	defer wg.Done()
	key := metric.GetDefinitionKey(event)
	metricValue, ok := wapi.Storage.Get(key)
	eventValue, _ := strconv.Atoi(fmt.Sprintf("%v", event.Data[metric.Aggregation]))
	if inputEvent {
		if !ok {
			metricValue = storage.MetricItem{Key: key, Value: (1<<bits.UintSize)/2 - 1}
		}
		if eventValue < metricValue.Value {
			metricValue.Value = eventValue
			wapi.Storage.Put(key, metricValue)
		}
	} else if eventValue == metricValue.Value { // check if the event leaving the window is the current min
		partialValue := (1<<bits.UintSize)/2 - 1
		for _, item := range wapi.deltaSlice.Items {
			itemValue, _ := strconv.Atoi(fmt.Sprintf("%v", item.Data[metric.Aggregation]))
			if itemValue < partialValue && itemValue != metricValue.Value {
				partialValue = itemValue
			}
		}
		metricValue.Value = partialValue
		wapi.Storage.Put(key, metricValue)
	}
}

func (wapi *WorkerAPI) doCount(metric storage.MetricDefinition, event storage.Event, inputEvent bool, wg *sync.WaitGroup) {
	key := metric.GetDefinitionKey(event)
	defer wg.Done()
	metricValue, ok := wapi.Storage.Get(key)
	if inputEvent {
		if !ok {
			metricValue = storage.MetricItem{Key: key, Value: 0}
		}
		metricValue.Value++
	} else {
		metricValue.Value--
	}
	//saving calculated metric value
	wapi.Storage.Put(key, metricValue)
}

func (wapi *WorkerAPI) doSum(metric storage.MetricDefinition, event storage.Event, inputEvent bool, wg *sync.WaitGroup) {
	key := metric.GetDefinitionKey(event)
	defer wg.Done()
	metricValue, ok := wapi.Storage.Get(key)
	aggregationValue, _ := strconv.Atoi(fmt.Sprintf("%v", event.Data[metric.Aggregation]))
	if inputEvent {
		if !ok {
			metricValue = storage.MetricItem{Key: key, Value: 0}
		}

		metricValue.Value = metricValue.Value + aggregationValue
	} else {
		metricValue.Value = metricValue.Value - aggregationValue
	}
	//saving calculated metric value
	wapi.Storage.Put(key, metricValue)
}
