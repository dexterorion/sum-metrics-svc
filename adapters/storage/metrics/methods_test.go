package metrics_storage

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/dexterorion/sum-metrics-svc/internal/core/domain"
	"github.com/dexterorion/sum-metrics-svc/internal/core/ports"
	"github.com/stretchr/testify/require"
)

var (
	ctx                 context.Context
	metricsInMemStorage ports.MetricsRepository
)

func TestMain(m *testing.M) {
	ctx = context.Background()

	// for the sake of testability, will add a small time to clean
	metricsInMemStorage = NewMetricsInMemStorage(1000 * time.Millisecond)

	code := m.Run()

	os.Exit(code)
}

func TestMetricsInMemStorageKeepAll(t *testing.T) {
	ticker := time.NewTicker(10 * time.Millisecond)
	tickCounter := 0

	for {
		<-ticker.C

		if tickCounter == 10 {
			break
		}

		metricsInMemStorage.AddMetric(ctx, &domain.Metric{
			Key:   "test",
			Value: 10,
		})

		if tickCounter%2 == 0 {
			metricsInMemStorage.AddMetric(ctx, &domain.Metric{
				Key:   "test2",
				Value: 10,
			})
		}

		if tickCounter%5 == 0 {
			metricsInMemStorage.AddMetric(ctx, &domain.Metric{
				Key:   "test5",
				Value: 10,
			})
		}

		if tickCounter == 9 {
			metricsInMemStorage.AddMetric(ctx, &domain.Metric{
				Key:   "test9",
				Value: 10,
			})
		}

		tickCounter++
	}

	test, err := metricsInMemStorage.GetMetricSum(ctx, "test")
	require.NoError(t, err)
	require.Equal(t, "test", test.Key)
	require.Equal(t, int64(100), test.Value)

	test2, err := metricsInMemStorage.GetMetricSum(ctx, "test2")
	require.NoError(t, err)
	require.Equal(t, "test2", test2.Key)
	require.Equal(t, int64(50), test2.Value)

	test5, err := metricsInMemStorage.GetMetricSum(ctx, "test5")
	require.NoError(t, err)
	require.Equal(t, "test5", test5.Key)
	require.Equal(t, int64(20), test5.Value)

	test9, err := metricsInMemStorage.GetMetricSum(ctx, "test9")
	require.NoError(t, err)
	require.Equal(t, "test9", test9.Key)
	require.Equal(t, int64(10), test9.Value)

	notest, err := metricsInMemStorage.GetMetricSum(ctx, "notest")
	require.NoError(t, err)
	require.Equal(t, "notest", notest.Key)
	require.Equal(t, int64(0), notest.Value)
}

func TestMetricsInMemStorageCleanSome(t *testing.T) {
	ticker := time.NewTicker(250 * time.Millisecond)
	tickCounter := 0

	for {
		<-ticker.C

		if tickCounter == 4 {
			break
		}

		metricsInMemStorage.AddMetric(ctx, &domain.Metric{
			Key:   "keepsome",
			Value: 10,
		})

		tickCounter++
	}

	keepsome, err := metricsInMemStorage.GetMetricSum(ctx, "keepsome")
	require.NoError(t, err)
	require.Equal(t, "keepsome", keepsome.Key)
	require.LessOrEqual(t, keepsome.Value, int64(40)) // sometimes cleaner runs faster so value might be 0
}
