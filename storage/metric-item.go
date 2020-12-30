package storage

// MetricItem adsadasd
type MetricItem struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
	Delta int    `json:"delta"`
}

//GetItens asdasd
func (api *StorageAPI) GetItens() map[string]MetricItem {
	return api.metricItens.GetMap()
}
