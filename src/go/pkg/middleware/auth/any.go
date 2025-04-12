package auth

import "net/http"

type AnyOrNoAuth struct{}

func (fa AnyOrNoAuth) Authenticate(_ *http.Request) (bool, error) {
	return true, nil
}
