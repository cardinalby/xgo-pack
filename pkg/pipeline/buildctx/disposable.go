package buildctx

import (
	"errors"
	"slices"
	"sync"
)

type Disposable interface {
	Dispose() error
}

type Disposables struct {
	items []Disposable
	mu    sync.Mutex
}

func NewDisposables() *Disposables {
	return &Disposables{}
}

func (ds *Disposables) Add(d Disposable) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.items = append(ds.items, d)
}

func (ds *Disposables) AddIfDisposable(obj any) {
	if d, ok := obj.(Disposable); ok {
		ds.Add(d)
	}
}

func (ds *Disposables) Dispose() error {
	var errs []error
	ds.mu.Lock()
	defer ds.mu.Unlock()
	var notDisposed []Disposable
	// iterate in reverse order to dispose in the opposite order of adding
	for i := len(ds.items) - 1; i >= 0; i-- {
		if e := ds.items[i].Dispose(); e != nil {
			errs = append(errs, e)
		} else {
			notDisposed = append(notDisposed, ds.items[i])
		}
	}
	slices.Reverse(notDisposed)
	ds.items = notDisposed

	return errors.Join(errs...)
}
