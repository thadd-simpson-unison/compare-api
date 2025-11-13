package main

import (
	"fmt"
	"net/http"

	"github.com/rs/cors"

	"compare-api/Routes"
	"compare-api/Swagger"
	"compare-api/Utilities"
)

func main() {
	Utilities.GlobalConfig = Utilities.GetConfig()

	c := Utilities.GlobalConfig
	port := ":" + c.Port
	baseUrl := c.Host + port

	corsHandler := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedOrigins: []string{baseUrl},
	})

	mux := Routes.AddRoutes()
	handler := corsHandler.Handler(mux)

	Swagger.Setup(c, baseUrl, mux)
	err := http.ListenAndServe(port, handler)
	if err != nil {
		fmt.Printf(err.Error())
	}
}
