package ports

import (
	"context"

	"github.com/dexterorion/sum-metrics-svc/internal/core/domain"
)

type MetricsUpdate interface {
	AddMetric(ctx context.Context, metric *domain.Metric) error
	GetMetricSum(ctx context.Context, metricKey string) (*domain.Metric, error)
}
