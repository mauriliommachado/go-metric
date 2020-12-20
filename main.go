package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/bits"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// Metric asdasd
type Metric struct {
	Target      string `json:"target"`
	Aggregation string `json:"aggregation"`
	MetricType  string `json:"metric_type"`
	Delta       int    `json:"delta"`
}

// MetricItem adsadasd
type MetricItem struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
	Delta int    `json:"delta"`
}

// WorkerAPI asdas
type WorkerAPI struct {
	input   chan map[string]interface{}
	results chan int
	quit    chan int
}

const (
	numberOfWorkers = 1
)

var metricstore []Metric
var metricitens map[string]MetricItem
var wapi WorkerAPI

func main() {
	fmt.Printf("hello, world\n")
	initMetrics()
	wapi.input = make(chan map[string]interface{}, 100)
	wapi.results = make(chan int, 10)
	wapi.quit = make(chan int)
	for i := 1; i <= numberOfWorkers; i++ {
		go worker(i)
	}
	initAPI()
}

func submitWork(data map[string]interface{}) {
	wapi.input <- data
}

func getKey(event map[string]interface{}, metric Metric) string {
	var target = fmt.Sprintf("%v", event[metric.Target])
	var sb strings.Builder
	sb.WriteString(metric.MetricType)
	if metric.Aggregation != "" {
		sb.WriteString("_" + metric.Aggregation)
	}
	if target != "<nil>" {
		sb.WriteString("_" + target)
	}
	return sb.String()
}

func initMetrics() {
	metricstore = make([]Metric, 0)
	metricitens = make(map[string]MetricItem, 0)
	var metric1 = Metric{Target: "client_id", Aggregation: "", MetricType: "count", Delta: 86400000}
	var metric2 = Metric{Target: "client_id", Aggregation: "amount", MetricType: "sum", Delta: 86400000}
	var metric3 = Metric{Target: "client_id", Aggregation: "amount", MetricType: "min", Delta: 86400000}
	var metric4 = Metric{Target: "client_id", Aggregation: "amount", MetricType: "max", Delta: 86400000}
	var metric5 = Metric{Target: "", Aggregation: "amount", MetricType: "min", Delta: 86400000}
	var metric6 = Metric{Target: "", Aggregation: "amount", MetricType: "max", Delta: 86400000}
	var metric7 = Metric{Target: "", Aggregation: "", MetricType: "count", Delta: 86400000}
	metricstore = append(metricstore, metric1)
	metricstore = append(metricstore, metric2)
	metricstore = append(metricstore, metric3)
	metricstore = append(metricstore, metric4)
	metricstore = append(metricstore, metric5)
	metricstore = append(metricstore, metric6)
	metricstore = append(metricstore, metric7)
	fmt.Println("Metric store initialized")
}

func worker(id int) {
	fmt.Println("Worker", id, " waiting for work")
	for {
		select {
		case event := <-wapi.input:
			fmt.Println(event, " data received on worker ", id)
			for _, metric := range metricstore {
				switch {
				case metric.MetricType == "count":
					doCount(metric, event)
				case metric.MetricType == "sum":
					doSum(metric, event)
				case metric.MetricType == "max":
					doMax(metric, event)
				case metric.MetricType == "min":
					doMin(metric, event)
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

func doMax(metric Metric, event map[string]interface{}) {
	key := getKey(event, metric)
	metricValue, ok := metricitens[key]
	if !ok {
		metricValue = MetricItem{Key: key, Value: (1 << bits.UintSize) / -2}
	}
	eventValue, _ := strconv.Atoi(fmt.Sprintf("%v", event[metric.Aggregation]))
	if eventValue > metricValue.Value {
		metricValue.Value = eventValue
		//saving calculated metric value
		metricitens[key] = metricValue
	}
}

func doMin(metric Metric, event map[string]interface{}) {
	key := getKey(event, metric)
	metricValue, ok := metricitens[key]
	if !ok {
		metricValue = MetricItem{Key: key, Value: (1<<bits.UintSize)/2 - 1}
	}
	eventValue, _ := strconv.Atoi(fmt.Sprintf("%v", event[metric.Aggregation]))
	if eventValue < metricValue.Value {
		metricValue.Value = eventValue
		//saving calculated metric value
		metricitens[key] = metricValue
	}
}

func doCount(metric Metric, event map[string]interface{}) {
	key := getKey(event, metric)
	metricValue, ok := metricitens[key]
	if !ok {
		metricValue = MetricItem{Key: key, Value: 0}
	}
	metricValue.Value++
	//saving calculated metric value
	metricitens[key] = metricValue
}

func doSum(metric Metric, event map[string]interface{}) {
	key := getKey(event, metric)
	metricValue, ok := metricitens[key]
	if !ok {
		metricValue = MetricItem{Key: key, Value: 0}
	}
	aggregationValue, _ := strconv.Atoi(fmt.Sprintf("%v", event[metric.Aggregation]))
	metricValue.Value = metricValue.Value + aggregationValue
	//saving calculated metric value
	metricitens[key] = metricValue
}

func initAPI() {
	router := httprouter.New()
	router.PUT("/metric", index)
	router.GET("/metric", get)
	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func get(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	j, _ := json.Marshal(metricitens)
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", j)
}

func index(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var data map[string]interface{}
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &data)
	submitWork(data)
	writeOKResponse(w)
}

func writeOKResponse(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// Writes the error response as a Standard API JSON response with a response code
func writeErrorResponse(w http.ResponseWriter, errorCode int, errorMsg string) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(errorCode)
}
