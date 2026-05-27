package httpx

import (
	"net/http"

	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/valueobject"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := r.Header.Get("X-User-ID")
		if raw == "" {
			WriteError(w, myerror.NewUnauthorizedError().SetMessage("X-User-ID header is required"))
			return
		}
		id, err := valueobject.NewPrimaryIdFromString[valueobject.UserID](raw)
		if err != nil {
			WriteError(w, myerror.NewUnauthorizedError().SetMessage("invalid X-User-ID header"))
			return
		}

		ctx := WithUserID(r.Context(), *id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
