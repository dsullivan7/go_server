package middlewares

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
)

func (m *Middlewares) Logger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()

			defer func() {
				meta := map[string]interface{}{
					"proto":        r.Proto,
					"method":       r.Method,
					"path":         r.URL.Path,
					"query":        r.URL.Query(),
					"responseTime": time.Since(t1),
					"status":       ww.Status(),
					"size":         ww.BytesWritten(),
					"reqId":        middleware.GetReqID(r.Context()),
				}

				m.logger.InfoWithMeta("Response", meta)
			}()

			next.ServeHTTP(ww, r)
		})
	}
}
