package teststeps

import (
	"github.com/vandad1901/p3s/apps/auth/internal/identity"
	"github.com/vandad1901/p3s/apps/auth/internal/scenario"
	"github.com/vandad1901/p3s/packages/go/gen/protobuf/auth/authnpb/v1"
	"github.com/vandad1901/p3s/packages/go/godogutil"
)

func resolveUser(s *godogutil.SharedData, row map[string]string) *identity.User {
	u := new(identity.User)

	u.Username = godogutil.ResolveString(s, row, "Username")
	u.Email = godogutil.ResolveString(s, row, "Email")

	return u
}

func resolveRefreshJWTRequest(s *scenario.Scenario, row map[string]string) *authnpb.RefreshJWTRequest {
	res := new(authnpb.RefreshJWTRequest)

	res.UserId = godogutil.ResolveInt64(s.SharedData, row, "UserID")
	res.SessionId = godogutil.ResolveInt64(s.SharedData, row, "SessionID")
	res.RefreshToken = godogutil.ResolveString(s.SharedData, row, "RefreshToken")

	return res
}
