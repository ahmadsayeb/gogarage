package mux

import (
	"os"

	"github.com/ardanlabs/service/apis/services/sales/route/sys/checkapi"
	"github.com/ardanlabs/service/foundation/web"
)

func WebAPI(shutdown chan os.Signal) *web.App {
	app := web.NewApp(shutdown)
	checkapi.Route(app)
	return app
}
