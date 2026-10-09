package authn

import "github.com/vandad1901/p3s/packages/go/apperror"

var (
	errClosedRegistration = apperror.Unauthenticated("authn.closedRegistration")
	errInvalidAuthn       = apperror.Unauthenticated("authn.invalidAuth")
)
