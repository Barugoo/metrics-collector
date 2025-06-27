package store

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/skdiver33/metrics-collector/models"
)

type MemStorage struct {
	storage map[string]models.Metrics
	mu      *sync.RWMutex
	ts      int64 // timestamp for metrics, used for testing purposes
}

type Storage interface {
	AddMetrics(metricsName string, metricsValue models.Metrics) error
	UpdateMetrics(metricsName string, metricsValue models.Metrics) error
	GetMetrics(metricsName string) (models.Metrics, error)
	GetAllMetricsNames() ([]string, error)
}

var ErrNotFound = errors.New("metrics not found")

func NewMemStorage(ts int64) Storage {
	return &MemStorage{
		storage: make(map[string]models.Metrics),
		mu:      &sync.RWMutex{},
		ts:      ts,
	}
}

func (inMemmory *MemStorage) AddMetrics(metricsName string, metricsValue models.Metrics) error {
	inMemmory.mu.Lock()
	defer inMemmory.mu.Unlock()

	_, ok := inMemmory.storage[metricsName]
	if ok {
		message := fmt.Sprintf("Metrics with name %s already exists", metricsName)
		return errors.New(message)
	}

	inMemmory.storage[metricsName] = metricsValue
	return nil

}

func (inMemmory *MemStorage) GetMetrics(metricsName string) (models.Metrics, error) {
	inMemmory.mu.Lock()
	defer inMemmory.mu.Unlock()

	metrics, ok := inMemmory.storage[metricsName]
	if !ok {
		metrics = models.Metrics{}
		return metrics, fmt.Errorf("%s: %w", metricsName, ErrNotFound)
	}
	return metrics, nil
}

func Safe[T any](v *T) T {
	if v == nil {
		var zero T
		return zero
	}
	return *v
}

func (inMemmory *MemStorage) UpdateMetrics(metricsName string, metricsValue models.Metrics) error {
	inMemmory.mu.Lock()
	defer inMemmory.mu.Unlock()

	_, ok := inMemmory.storage[metricsName]
	if !ok {
		message := fmt.Sprintf("Metrics with name %s not found", metricsName)
		return errors.New(message)
	}

	log.Printf("%d Update metrics %s, type: %s, value: %.6f, delta: %d",
		inMemmory.ts, metricsName, metricsValue.MType, Safe(metricsValue.Value), Safe(metricsValue.Delta))
	inMemmory.storage[metricsName] = metricsValue
	return nil
}

func (inMemmory *MemStorage) GetAllMetricsNames() ([]string, error) {
	inMemmory.mu.Lock()
	defer inMemmory.mu.Unlock()

	var allMetricsName []string
	for name := range inMemmory.storage {
		allMetricsName = append(allMetricsName, name)
	}
	return allMetricsName, nil

}
