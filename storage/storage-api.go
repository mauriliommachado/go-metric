package storage

// API asdas
type API struct {
	metricDefinitions []MetricDefinition
	metricItens       ValueHashtable
}

//New todo
func New() *API {
	var strAPI API
	strAPI.metricDefinitions = make([]MetricDefinition, 0)
	strAPI.metricItens = ValueHashtable{}
	return &strAPI
}

//GetMetricDefinitions todo
func (api *API) GetMetricDefinitions() []MetricDefinition {
	return api.metricDefinitions
}

// Put item with value v and key k into the hashtable
func (api *API) Put(k string, v MetricItem) {
	api.metricItens.Put(k, v)
}

// Remove item with key k from hashtable
func (api *API) Remove(k string) {
	api.metricItens.Remove(k)
}

// Get item with key k from the hashtable
func (api *API) Get(k string) (MetricItem, bool) {
	return api.metricItens.Get(k)
}

// GetMap from table
func (api *API) GetMap() map[string]MetricItem {
	return api.metricItens.GetMap()
}
