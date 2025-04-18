package checkapi

import "net/http"

func Route(mux *http.ServeMux) {
	mux.HandleFunc("/liveness", liveness)
	mux.HandleFunc("/readiness", readiness)

}
