package main

import (
	"flag"
	"net/http"
	"time"

	web_handler "github.com/dexterorion/sum-metrics-svc/adapters/api/web"
	metrics_handler "github.com/dexterorion/sum-metrics-svc/adapters/api/web/handlers/metrics"
	metrics_storage "github.com/dexterorion/sum-metrics-svc/adapters/storage/metrics"
	"github.com/dexterorion/sum-metrics-svc/internal/core/usecases"
	"github.com/go-openapi/spec"

	"github.com/dexterorion/sum-metrics-svc/pkg/logging"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	restful "github.com/emicklei/go-restful/v3"
)

var (
	log = logging.Init("metrics_service")

	withswagger bool
	swaggerdir  string
	binding     string
)

func init() {
	flag.BoolVar(&withswagger, "withswagger", true, "Creates swagger and installs service on localhost:8080/apidocs.json")
	flag.StringVar(&swaggerdir, "swaggerdir", "", "Sets projects base path where swagger should be")
	flag.StringVar(&binding, "binding", ":8080", "Web service port")
	flag.Parse()
}

func main() {

	log.Infow("Bootstraping with flags",
		"withswagger", withswagger,
		"swaggerdir", swaggerdir,
		"binding", binding,
	)

	handler, err := web_handler.NewWebHandler(func(ws *restful.WebService) error {
		// metrics handler injection
		metricsStorage := metrics_storage.NewMetricsInMemStorage(time.Hour)
		metricsUpdateUC := usecases.NewMetricsUpdate(metricsStorage)
		metrics_handler.NewMetricsHandler(ws, metricsUpdateUC)

		return nil
	})

	if withswagger {
		config := restfulspec.Config{
			WebServices:                   restful.RegisteredWebServices(), // you control what services are visible
			APIPath:                       "/apidocs.json",
			PostBuildSwaggerObjectHandler: enrichSwaggerObject}
		restful.DefaultContainer.Add(restfulspec.NewOpenAPIService(config))

		// Optionally, you can install the Swagger Service which provides a nice Web UI on your REST API
		// You need to download the Swagger HTML5 assets and change the FilePath location in the config below.
		// Open http://localhost:8080/apidocs/?url=http://localhost:8080/apidocs.json
		http.Handle("/apidocs/", http.StripPrefix("/apidocs/", http.FileServer(http.Dir(swaggerdir))))

		// Optionally, you may need to enable CORS for the UI to work.
		cors := restful.CrossOriginResourceSharing{
			AllowedHeaders: []string{"Content-Type", "Accept"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
			CookiesAllowed: false,
			Container:      restful.DefaultContainer}
		restful.DefaultContainer.Filter(cors.Filter)
	}

	if err != nil {
		log.Fatalw("error starting web hander", "error", err)
	}

	handler.StartBlocking(binding)
}

func enrichSwaggerObject(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "MetricService",
			Description: "Resource for managing metric",
			Contact: &spec.ContactInfo{
				ContactInfoProps: spec.ContactInfoProps{
					Name:  "john",
					Email: "john@doe.rp",
					URL:   "http://johndoe.org",
				},
			},
			License: &spec.License{
				LicenseProps: spec.LicenseProps{
					Name: "MIT",
					URL:  "http://mit.org",
				},
			},
			Version: "1.0.0",
		},
	}
	swo.Tags = []spec.Tag{
		{
			TagProps: spec.TagProps{
				Name:        "metrics",
				Description: "Managing metrics",
			},
		},
	}
}
