package checkapi

import (
	"github.com/ardanlabs/service/foundation/web"
)

func Route(app *web.App) {
	app.HandleFuncNoMiddleware("/liveness", liveness)
	app.HandleFuncNoMiddleware("/readiness", readiness)
	app.HandleFunc("/testerror", testerror)
	app.HandleFunc("/testpanic", testpanic)
}
