package middlewares

import (
	"go_server/internal/logger"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
)

func Logger(l logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()

			defer func() {
				meta := map[string]interface{}{
					"proto":  r.Proto,
					"path":   r.URL.Path,
					"query":  r.URL.Query(),
					"lat":    time.Since(t1),
					"status": ww.Status(),
					"size":   ww.BytesWritten(),
					"reqId":  middleware.GetReqID(r.Context()),
				}

				l.InfoWithMeta("Response", meta)
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
