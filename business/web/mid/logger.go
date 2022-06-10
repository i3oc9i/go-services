package mid

import (
	"context"
	"net/http"
	"time"

	"github.com/i3oc9i/go-service/foundation/web"
	"go.uber.org/zap"
)

// Logger ...
func Logger(log *zap.SugaredLogger) web.Middleware {

	m := func(handler web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			traceID := "00000000-0000-0000-0000-000000000000"
			statusCode := http.StatusOK
			now := time.Now()

			log.Infow("request started", "traceID", traceID, "method", r.Method, "path", r.URL.Path,
				"remoteaddr", r.RemoteAddr)

			err := handler(ctx, w, r)

			log.Infow("request complted", "traceID", traceID, "method", r.Method, "path", r.URL.Path,
				"remoteaddr", r.RemoteAddr, "statuscode", statusCode, "since", time.Since(now))

			return err
		}

		return h
	}

	return m
}
