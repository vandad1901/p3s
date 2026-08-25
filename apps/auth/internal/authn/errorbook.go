package authn

import "github.com/vandad1901/p3s/packages/go/apperror"

var (
	errInvalidAuthn = apperror.Unauthenticated("authn.invalidAuth")
)
