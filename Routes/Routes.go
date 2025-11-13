package Routes

import (
	"compare-api/Routes/Handlers"
	"net/http"
)

func AddRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/compare", Handlers.CompareHandler)

	return mux
}
