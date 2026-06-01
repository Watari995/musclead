package authhandler

import (
	"net/http"
	"strings"

	sessiondomain "github.com/Watari995/musclead/internal/auth/internal/domain"
	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/shared/httpx"
)

const jwtPrefix = "Bearer "

func NewJWTMiddleware(signer sessiondomain.TokenSigner) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Authorization ヘッダを取得する
			authHeader := r.Header.Get("Authorization")
			// "Bearer " で始まるかtrimをしてチェック
			if !strings.HasPrefix(authHeader, jwtPrefix) {
				httpx.WriteError(w, myerror.NewUnauthorizedError())
				return
			}
			tokenStr := strings.TrimPrefix(authHeader, jwtPrefix)
			// jwt検証 -> userIDを取得する
			userID, err := signer.VerifyAccessToken(tokenStr)
			if err != nil {
				httpx.WriteError(w, myerror.NewUnauthorizedError().SetMessage("failed to verify token"))
				return
			}

			ctx := httpx.WithUserID(r.Context(), userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
