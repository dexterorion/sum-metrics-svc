package web_handler

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dexterorion/sum-metrics-svc/pkg/logging"
	"github.com/emicklei/go-restful/v3"
)

var (
	log = logging.Init("web_handler")
)

type WebHandler struct {
	ws *restful.WebService
}

func NewWebHandler(injectHandlers func(ws *restful.WebService) error) (*WebHandler, error) {
	ws := new(restful.WebService)
	ws = ws.Path("/api")

	if injectHandlers != nil {
		err := injectHandlers(ws)
		if err != nil {
			return nil, err
		}
	}

	restful.DefaultContainer.Add(ws)

	return &WebHandler{
		ws: ws,
	}, nil
}

func (h *WebHandler) GetWs() *restful.WebService {
	return h.ws
}

func (h *WebHandler) StartBlocking(port string) {
	go func() {
		log.Infow("Listening on port " + port)
		log.Fatal(http.ListenAndServe(port, nil))
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool, 1)
	go func() {
		sig := <-signals
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	<-done
	os.Exit(0)
}
