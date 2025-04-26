package web

import (
	"context"
	"errors"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/google/uuid"
)

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App struct.
type App struct {
	// embedded value
	*http.ServeMux
	shutdown chan os.Signal
	mw       []MidHandler
}

func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

func NewApp(shutdown chan os.Signal, mw ...MidHandler) *App {
	return &App{
		ServeMux: http.NewServeMux(),
		shutdown: shutdown,
		mw:       mw,
	}
}

func (a *App) HandleFunc(pattern string, handler Handler, mw ...MidHandler) {
	handler = wrapMiddleware(mw, handler)
	handler = wrapMiddleware(a.mw, handler)
	h := func(w http.ResponseWriter, r *http.Request) {
		v := Values{
			TraceID: uuid.NewString(),
			Now:     time.Now().UTC(),
		}
		ctx := setValues(r.Context(), &v)
		//PUT CODE HERE
		// Logging Direct NO! foundation package shouldn't have logging.
		// thats why we have the middleware
		if err := handler(ctx, w, r); err != nil {
			// fmt.Println("Error handling request:", err)
			if validateError(err) {
				a.SignalShutdown()
				return
			}
		}
	}
	a.ServeMux.HandleFunc(pattern, h)
}

func validateError(err error) bool {
	switch {
	case errors.Is(err, syscall.EPIPE):
		return false
	case errors.Is(err, syscall.ECONNRESET):
		return false

	}
	return true
}
