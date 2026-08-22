package post_test

import (
	"context"
	"testing"

	"github.com/cucumber/godog"
	"github.com/vandad1901/p3s/apps/api/internal/app"
	"github.com/vandad1901/p3s/apps/api/internal/post/teststeps"
	"github.com/vandad1901/p3s/apps/api/internal/scenario"
	"github.com/vandad1901/p3s/packages/go/godogutil"
	"github.com/vandad1901/p3s/packages/go/usercontext"
)

func TestFeatures(t *testing.T) { //nolint:paralleltest
	suite := godog.TestSuite{
		Name: "features",
		ScenarioInitializer: func(sc *godog.ScenarioContext) {
			InitializeScenario(t, sc)
		},
		Options: &godog.Options{
			Paths:          []string{"features"},
			TestingT:       t,
			Format:         "pretty",
			DefaultContext: usercontext.CtxWithUser(context.Background(), 1),
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
		A:          app.Boot(),
	}

	godogutil.DefineStep(s.SharedData, ctx, "post")
	godogutil.AddStep(ctx, "post blocks", godogutil.AddToHeaderStep(s.SharedData, teststeps.AddedPostBlocksKey))
	godogutil.EditStep(ctx, "post blocks", godogutil.AddToHeaderStep(s.SharedData, teststeps.EditedPostBlocksKey))
	godogutil.RemoveStep(ctx, "post blocks", godogutil.AddToHeaderStep(s.SharedData, teststeps.RemovedPostBlocksKey))
	godogutil.CreateStep(ctx, "post", teststeps.CreateStep(s))
	godogutil.GetStep(ctx, "post", teststeps.GetStep(s))
	godogutil.MustBeDeleted(ctx, "post", teststeps.MustBeDeletedStep(s))
	godogutil.UpdateStep(ctx, "post", teststeps.UpdateStep(s))
	godogutil.DeleteStep(ctx, "post", teststeps.DeleteStep(s))

	godogutil.ReturnedErrorStep(ctx, s.SharedData)
}
