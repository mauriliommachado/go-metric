package engine

import (
	"testing"

	"github.com/mauriliommachado/go-metric/storage"
)

func TestCleanUpSlice(t *testing.T) {
	wapi := New(1)
	data := map[string]interface{}{}
	ok := wapi.deltaSlice.Push(storage.Event{1, data})
	if !ok {
		t.Errorf("Push failed")
	}
	t.Log(wapi.deltaSlice.Items)
	ok = wapi.deltaSlice.Push(storage.Event{2, data})
	if !ok {
		t.Errorf("Push failed")
	}
	t.Log(wapi.deltaSlice.Items)
	ok = wapi.deltaSlice.Push(storage.Event{3, data})
	if !ok {
		t.Errorf("Push failed")
	}
	t.Log(wapi.deltaSlice.Items)
	ok = wapi.deltaSlice.Push(storage.Event{5, data})
	if !ok {
		t.Errorf("Push failed")
	}
	wapi.deltaSlice.cleanUpSlice()
	t.Log(wapi.deltaSlice.Items)
}
