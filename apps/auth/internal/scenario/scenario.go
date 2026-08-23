package scenario

import (
	"github.com/vandad1901/p3s/apps/auth/internal/app"
	"github.com/vandad1901/p3s/packages/go/godogutil"
)

type Scenario struct {
	*godogutil.SharedData

	A *app.App
}
