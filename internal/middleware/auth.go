package middleware

import (
	"fmt"
	"net/http"

	"github.com/nabind47/sensor_api/internal/config"
	"github.com/nabind47/sensor_api/internal/util"
)

type Auth struct {
	cfg *config.AuthConfig
}

func NewAuth(cfg *config.AuthConfig) *Auth {
	return &Auth{cfg: cfg}
}

func (a *Auth) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("x-authorization-key")

		fmt.Printf("token %v", token)

		if token == "" {
			util.WriteError(w, http.StatusForbidden, "forbidden")
			return
		}

		ok := util.ValidateHash(a.cfg.ClientID, a.cfg.ClientSecret, a.cfg.ClientExpiry, token)
		if !ok {
			util.WriteError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		next.ServeHTTP(w, r)
	})
}
