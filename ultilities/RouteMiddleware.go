package ultilities

import (
	ProviderJWT "WeddingBackEnd/ultilities/provider/jwt"
	"context"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
)

func Midlleware(JWTKEY string) func(next httprouter.Handle) httprouter.Handle {
	return func(next httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			var claims ProviderJWT.Claims
			header := r.Header.Get("Authorization")
			tokenString := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(header), "Bearer"))

			token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(JWTKEY), nil
			})

			if err != nil || token == nil || !token.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				w.Header().Set("Content-Type", "Token invalid")
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), "user.email", claims.Subject))
			next(w, r, p)
		}
	}
}
