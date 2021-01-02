package engine

import (
	"sort"
	"sync"

	"github.com/mauriliommachado/go-metric/storage"
)

//newDeltaSlice a new DeltaSlice
func (wapi *WorkerAPI) newDeltaSlice(delta int64) {
	var ds DeltaSlice
	ds.Delta = delta
	ds.Items = []storage.Event{}
	ds.clock = 0
	ds.api = wapi
	wapi.deltaSlice = &ds
}

//DeltaSlice todo
type DeltaSlice struct {
	Items []storage.Event
	Delta int64
	clock int64
	api   *WorkerAPI
}

func (s *DeltaSlice) updateClock(event storage.Event) {
	if event.Timestamp > s.clock {
		s.clock = event.Timestamp
	}
	s.cleanUpSlice()
}

func (s *DeltaSlice) cleanUpSlice() {
	var indexToRemove int
	indexToRemove = -1
	for i, event := range s.Items {
		if !s.EventBelongs(event) {
			var waitgroup sync.WaitGroup
			s.api.do(event, false)
			waitgroup.Wait()
			indexToRemove = i
		} else {
			break
		}
	}
	if indexToRemove > -1 {
		s.Items = s.Items[indexToRemove+1:]
	}
}

//Push an item in order
func (s *DeltaSlice) Push(event storage.Event) bool {
	s.updateClock(event)
	if !s.EventBelongs(event) {
		return false
	}
	index := sort.Search(len(s.Items),
		func(i int) bool { return s.Items[i].Timestamp > event.Timestamp })
	s.Items = append(s.Items, event)
	copy(s.Items[index+1:], s.Items[index:])
	s.Items[index] = event
	return true
}

//EventBelongs to the delta defined on the slice
func (s *DeltaSlice) EventBelongs(event storage.Event) bool {
	if s.IsEmpty() {
		return true
	}
	//fmt.Println("event with timestamp ", event.Timestamp, " with clock ", s.clock, " with delta ", s.Delta, " is ", event.Timestamp > (s.clock-s.Delta))
	return event.Timestamp > (s.clock - s.Delta)
}

//Pop an item in back
func (s *DeltaSlice) Pop() storage.Event {
	i := len(s.Items) - 1
	defer func() {
		s.Items = append(s.Items[:i], s.Items[i+1:]...)
	}()
	return s.Items[i]
}

//IsEmpty check if the deque is empty and his values
func (s *DeltaSlice) IsEmpty() bool {
	if len(s.Items) == 0 {
		return true
	}
	return false
}
