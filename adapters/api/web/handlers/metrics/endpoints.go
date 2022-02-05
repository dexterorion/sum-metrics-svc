package metrics_handler

import (
	"net/http"

	"github.com/dexterorion/sum-metrics-svc/adapters/api/shared"
	"github.com/dexterorion/sum-metrics-svc/internal/core/ports"
	"github.com/dexterorion/sum-metrics-svc/pkg/logging"
	restful "github.com/emicklei/go-restful/v3"
	"go.uber.org/zap"
)

type MetricsHandler struct {
	log           *zap.SugaredLogger
	metricsUpdate ports.MetricsUpdate
}

func NewMetricsHandler(ws *restful.WebService, metricsUpdate ports.MetricsUpdate) *MetricsHandler {
	handler := &MetricsHandler{
		log:           logging.Init("metrics_handler"),
		metricsUpdate: metricsUpdate,
	}

	ws.Route(
		shared.DefDefaultResponse(
			ws.POST("/metric/{key}").
				To((handler.PostMetric)).
				Param(ws.PathParameter("key", "Metric key").DataType("string")).
				Consumes(restful.MIME_JSON).
				Produces(restful.MIME_JSON).
				Reads(&NewMetricRequest{}, "").
				Returns(200, "Metric added", &shared.EmptyBody{})))

	ws.Route(
		shared.DefDefaultResponse(
			ws.GET("/metric/{key}/sum").
				To((handler.GetMetricSum)).
				Param(ws.PathParameter("key", "Metric key").DataType("string")).
				Consumes(restful.MIME_JSON).
				Produces(restful.MIME_JSON).
				Returns(200, "Metric sum result", &GetMetricSumResponse{})))

	return handler
}

func (h *MetricsHandler) GetMetricSum(req *restful.Request, resp *restful.Response) {
	key := req.PathParameter("key")

	metricSum, err := h.metricsUpdate.GetMetricSum(req.Request.Context(), key)
	if err != nil {
		h.log.Errorw("Error on getting metric sum", "key", key, "error", err)
		resp.WriteHeader(http.StatusInternalServerError)
		resp.WriteAsJson(&shared.ErrorResponse{Message: err.Error(), Code: http.StatusInternalServerError})
		return
	}

	response := &GetMetricSumResponse{}
	response.FromDomain(metricSum)
	resp.WriteAsJson(response)
}

func (h *MetricsHandler) PostMetric(req *restful.Request, resp *restful.Response) {
	var data *NewMetricRequest = &NewMetricRequest{}

	if err := req.ReadEntity(data); err != nil {
		h.log.Errorw("Error reading metric input", "error", err)
		resp.WriteHeader(http.StatusBadRequest)
		resp.WriteAsJson(&shared.ErrorResponse{Message: err.Error(), Code: http.StatusBadRequest})
		return
	}

	if err := data.Validate(); err != nil {
		h.log.Errorw("Error validating metric input", "error", err)
		resp.WriteHeader(http.StatusBadRequest)
		resp.WriteAsJson(&shared.ErrorResponse{Message: err.Error(), Code: http.StatusBadRequest})
		return
	}

	key := req.PathParameter("key")

	err := h.metricsUpdate.AddMetric(req.Request.Context(), data.ToDomain(key))

	if err != nil {
		h.log.Errorw("Error on adding new metric", "data", data, "error", err)
		resp.WriteHeader(http.StatusInternalServerError)
		resp.WriteAsJson(&shared.ErrorResponse{Message: err.Error(), Code: http.StatusInternalServerError})
		return
	}

	resp.WriteAsJson(&shared.EmptyBody{})
}
