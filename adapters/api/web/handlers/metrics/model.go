package metrics_handler

type NewMetricRequest struct {
	Value int64 `json:"value"`
}

type GetMetricSumResponse struct {
	Value int64 `json:"value"`
}
