package metrics_handler

import (
	"errors"

	"github.com/dexterorion/sum-metrics-svc/internal/core/domain"
)

type NewMetricRequest struct {
	Value int64 `json:"value"`
}

func (m *NewMetricRequest) Validate() error {
	if m.Value == 0 {
		return errors.New("`value` should be greater than 0")
	}

	return nil
}

func (m *NewMetricRequest) ToDomain(key string) *domain.Metric {
	return &domain.Metric{
		Key:   key,
		Value: m.Value,
	}
}

type GetMetricSumResponse struct {
	Value int64  `json:"value"`
	Key   string `json:"key"`
}

func (m *GetMetricSumResponse) FromDomain(data *domain.Metric) {
	if data == nil {
		return
	}

	m.Value = data.Value
	m.Key = data.Key
}
