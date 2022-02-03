package metrics_storage

import (
	"context"
	"sync"
	"time"

	"github.com/dexterorion/sum-metrics-svc/internal/core/domain"
	"github.com/dexterorion/sum-metrics-svc/internal/core/ports"
)

const (
	DefaultCleanTime = time.Hour
)

type entry struct {
	Date  time.Time
	Value int64
}

type MetricsInMemStorage struct {
	mu        sync.Mutex
	data      map[string][]*entry
	cleantime time.Duration
}

func NewMetricsInMemStorage(cleantime time.Duration) ports.MetricsRepository {
	return &MetricsInMemStorage{
		mu:        sync.Mutex{},
		data:      map[string][]*entry{},
		cleantime: cleantime,
	}
}

func (mim *MetricsInMemStorage) AddMetric(ctx context.Context, metric *domain.Metric) error {
	mim.mu.Lock()
	defer mim.mu.Unlock()

	// clean before inserting
	mim.cleanMetrics(metric.Key)

	if mim.data[metric.Key] == nil || len(mim.data[metric.Key]) == 0 {
		mim.data[metric.Key] = []*entry{}
	}

	mim.data[metric.Key] = append(mim.data[metric.Key], &entry{
		Value: metric.Value,
		Date:  time.Now(),
	})

	return nil
}

func (mim *MetricsInMemStorage) GetMetricSum(ctx context.Context, metricKey string) (*domain.Metric, error) {
	mim.mu.Lock()
	defer mim.mu.Unlock()

	// clean before calculating (because it will last only valid entries)
	mim.cleanMetrics(metricKey)

	var metricSum int64 = 0

	for _, v := range mim.data[metricKey] {
		metricSum += v.Value
	}

	return &domain.Metric{
		Key:   metricKey,
		Value: metricSum,
	}, nil
}

func (mim *MetricsInMemStorage) cleanMetrics(metricKey string) {
	metricValues := mim.data[metricKey]
	totalValues := len(metricValues)

	if metricValues == nil || totalValues == 0 {
		return
	}

	cuttime := time.Now().Add(-mim.cleantime)
	indexRemovalPoint := 0

	for ; indexRemovalPoint < totalValues; indexRemovalPoint++ {
		v := metricValues[indexRemovalPoint]

		if v.Date.Before(cuttime) {
			break
		}
	}

	mim.data[metricKey] = mim.data[metricKey][:indexRemovalPoint]
}
