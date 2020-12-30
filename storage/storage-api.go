package storage

// StorageAPI asdas
type StorageAPI struct {
	metricDefinitions []MetricDefinition
	metricItens       ValueHashtable
}

//New todo
func New() *StorageAPI {
	var strAPI StorageAPI
	strAPI.metricDefinitions = make([]MetricDefinition, 0)
	strAPI.metricItens = ValueHashtable{}
	return &strAPI
}

//GetMetricDefinitions todo
func (api *StorageAPI) GetMetricDefinitions() []MetricDefinition {
	return api.metricDefinitions
}

// Put item with value v and key k into the hashtable
func (api *StorageAPI) Put(k string, v MetricItem) {
	api.metricItens.Put(k, v)
}

// Remove item with key k from hashtable
func (api *StorageAPI) Remove(k string) {
	api.metricItens.Remove(k)
}

// Get item with key k from the hashtable
func (api *StorageAPI) Get(k string) (MetricItem, bool) {
	return api.metricItens.Get(k)
}

// GetMap from table
func (api *StorageAPI) GetMap() map[string]MetricItem {
	return api.metricItens.GetMap()
}
