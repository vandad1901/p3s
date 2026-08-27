package authn_test

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"

	"github.com/cucumber/godog"
	"github.com/vandad1901/p3s/apps/auth/internal/app"
	"github.com/vandad1901/p3s/apps/auth/internal/authn/teststeps"
	"github.com/vandad1901/p3s/apps/auth/internal/config"
	"github.com/vandad1901/p3s/apps/auth/internal/scenario"
	"github.com/vandad1901/p3s/packages/go/godogutil"
)

func TestFeatures(t *testing.T) { //nolint:paralleltest
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		"x-forwarded-for", "127.0.0.1",
		"user-agent", "gherkin-test",
	))

	suite := godog.TestSuite{
		Name: "features",
		ScenarioInitializer: func(sc *godog.ScenarioContext) {
			InitializeScenario(t, sc)
		},
		Options: &godog.Options{
			Paths:          []string{"features"},
			TestingT:       t,
			Format:         "pretty",
			DefaultContext: ctx,
		},
	}

	if status := suite.Run(); status != 0 {
		t.Fail()
	}
}

func InitializeScenario(t *testing.T, ctx *godog.ScenarioContext) {
	t.Helper()

	s := &scenario.Scenario{
		SharedData: godogutil.InitBaseData(t),
		A:          app.MustBoot(config.LoadConfig()),
	}

	ctx.Step(`^user registers with the following data( expecting error)?$`, teststeps.RegisterStep(s))
	ctx.Step(`^user should receive valid JWT with the following data$`, teststeps.AssertJWTStep(s))
	ctx.Step(`^user logs in with the following data( expecting error)?$`, teststeps.LoginStep(s))
	ctx.Step(`^user refreshes jwt with the following data( expecting error)?$`, teststeps.RefreshStep(s))

	godogutil.ReturnedErrorStep(ctx, s.SharedData)
}
