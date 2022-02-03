package usecases

import (
	"context"

	"github.com/dexterorion/sum-metrics-svc/internal/core/domain"
	"github.com/dexterorion/sum-metrics-svc/internal/core/ports"
	"github.com/dexterorion/sum-metrics-svc/pkg/logging"
	"go.uber.org/zap"
)

func NewMetricsUpdate(metricsStorage ports.MetricsRepository) *MetricsUpdateImpl {
	return &MetricsUpdateImpl{
		metricsStorage: metricsStorage,
		log:            logging.Init("metrics_update_uc"),
	}
}

type MetricsUpdateImpl struct {
	metricsStorage ports.MetricsRepository
	log            *zap.SugaredLogger
}

func (mu *MetricsUpdateImpl) AddMetric(ctx context.Context, metric *domain.Metric) error {
	mu.log.Infow("Adding a new metric", "key", metric.Key, "value", metric.Value)

	err := mu.metricsStorage.AddMetric(ctx, metric)
	if err != nil {
		mu.log.Infow("Error adding a new metric on storage", "key", metric.Key, "value", metric.Value, "err", err)
		return err
	}

	return nil
}

func (mu *MetricsUpdateImpl) GetMetricSum(ctx context.Context, metricKey string) (*domain.Metric, error) {
	mu.log.Infow("Getting metric sum", "key", metricKey)

	sum, err := mu.metricsStorage.GetMetricSum(ctx, metricKey)
	if err != nil {
		mu.log.Infow("Error getting metric sum from storage", "key", metricKey, "err", err)
		return nil, err
	}

	return sum, nil
}
