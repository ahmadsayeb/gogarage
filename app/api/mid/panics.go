package mid

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/ardanlabs/service/app/api/metrics"
)

func Panics(ctx context.Context, handler Handler) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			trace := debug.Stack()
			err = fmt.Errorf("PANIC [%v] TRACE[%s]", rec, string(trace))

			metrics.AddPanics(ctx)
		}
	}()
	return handler(ctx)
}
