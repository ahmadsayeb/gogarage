package checkapi

import (
	"github.com/ardanlabs/service/foundation/web"
)

func Route(app *web.App) {
	app.HandleFunc("/liveness", liveness)
	app.HandleFunc("/readiness", readiness)

}
