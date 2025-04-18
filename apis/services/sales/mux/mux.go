package mux

import (
	"net/http"

	"github.com/ardanlabs/service/apis/services/sales/route/sys/checkapi"
)

func WebAPI() *http.ServeMux {
	mux := http.NewServeMux()
	checkapi.Route(mux)
	return mux
}
