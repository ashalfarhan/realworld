package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/ashalfarhan/realworld/utils"
	"github.com/ashalfarhan/realworld/utils/logger"
	"github.com/google/uuid"
)

func InjectReqID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqID := strings.ReplaceAll(uuid.NewString(), "-", "")
		ctx := utils.CreateReqIDCtx(r.Context(), reqID)
		w.Header().Set("X-Request-ID", reqID)
		log := logger.GetCtx(ctx)
		log.Printf("Incoming %q request to %q", r.Method, r.URL.Path)
		log.Printf("QueyParams: %v", r.URL.Query())
		next.ServeHTTP(w, r.WithContext(ctx))
		end := time.Now()
		log.Printf("Request ends, took %dms", end.Sub(start).Milliseconds())
	})
}
