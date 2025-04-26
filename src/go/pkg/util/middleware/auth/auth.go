package auth

import "net/http"

type Authenticater interface {
	Authenticate(r *http.Request) (bool, error)
}

func WrapHttpAuth(auth Authenticater, handle http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ok, err := auth.Authenticate(r)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		handle(w, r)
	}
}
