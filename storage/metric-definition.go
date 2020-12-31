package storage

import (
	"fmt"
	"strings"
)

// MetricDefinition asdasd
type MetricDefinition struct {
	Target      string `json:"target"`
	Aggregation string `json:"aggregation"`
	MetricType  string `json:"metric_type"`
	Delta       int    `json:"delta"`
}

//GetDefinitionKey todo
func (metric MetricDefinition) GetDefinitionKey(event Event) string {
	var target = fmt.Sprintf("%v", event.Data[metric.Target])
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

//AddMetricDef todo
func (api *API) AddMetricDef(metricDefinition MetricDefinition) {
	api.metricDefinitions = append(api.metricDefinitions, metricDefinition)
}

//MockMetrics asdasd
func (api *API) MockMetrics() {
	api.AddMetricDef(MetricDefinition{Target: "client_id", Aggregation: "", MetricType: "count", Delta: 86400000})
	api.AddMetricDef(MetricDefinition{Target: "client_id", Aggregation: "amount", MetricType: "sum", Delta: 86400000})
	api.AddMetricDef(MetricDefinition{Target: "client_id", Aggregation: "amount", MetricType: "min", Delta: 86400000})
	api.AddMetricDef(MetricDefinition{Target: "client_id", Aggregation: "amount", MetricType: "max", Delta: 86400000})
	api.AddMetricDef(MetricDefinition{Target: "", Aggregation: "amount", MetricType: "min", Delta: 86400000})
	api.AddMetricDef(MetricDefinition{Target: "", Aggregation: "amount", MetricType: "max", Delta: 86400000})
	api.AddMetricDef(MetricDefinition{Target: "", Aggregation: "", MetricType: "count", Delta: 86400000})
}
