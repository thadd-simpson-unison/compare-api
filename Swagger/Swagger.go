package Swagger

import (
	"compare-api/Utilities"
	"net/http"

	"github.com/go-swagno/swagno"
	"github.com/go-swagno/swagno-http/swagger"
	"github.com/go-swagno/swagno/components/endpoint"
	"github.com/go-swagno/swagno/components/parameter"
)

func Setup(config Utilities.Config, hostname string, mux *http.ServeMux) {
	sw := swagno.New(swagno.Config{Title: config.ApiTitle, Version: config.ApiVersion, Host: hostname})

	sw.AddEndpoints(GetEndpoints())

	mux.HandleFunc("/swagger/", swagger.SwaggerHandler(sw.MustToJson()))
}

func GetEndpoints() []*endpoint.EndPoint {
	endpoints := []*endpoint.EndPoint{
		endpoint.New(
			endpoint.POST,
			"/compare",
			endpoint.WithSummary("Upload a CSV file."),
			endpoint.WithTags("External API Comparison"),
			endpoint.WithParams(
				parameter.FileParam("File"),
			),
		),
	}

	return endpoints
}
