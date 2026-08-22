//go:build test

package godogutil

import (
	"fmt"

	"github.com/cucumber/godog"
)

func CreateStep(ctx *godog.ScenarioContext, expr string, stepFunc any) {
	ctx.Step(fmt.Sprintf("^user creates %s with the following data( expecting error)?$", expr), stepFunc)
}

func GetStep(ctx *godog.ScenarioContext, expr string, stepFunc any) {
	ctx.Step(fmt.Sprintf("^user should be able to see %s with the following data( expecting error)?$", expr), stepFunc)
}

func MustBeDeleted(ctx *godog.ScenarioContext, expr string, stepFunc any) {
	ctx.Step(fmt.Sprintf("^user should not be able to see %s with the following data?$", expr), stepFunc)
}

func UpdateStep(ctx *godog.ScenarioContext, expr string, stepFunc any) {
	ctx.Step(fmt.Sprintf("^user updates %s with the following data( expecting error)?$", expr), stepFunc)
}

func DeleteStep(ctx *godog.ScenarioContext, expr string, stepFunc any) {
	ctx.Step(fmt.Sprintf("^user deletes %s with the following data( expecting error)?$", expr), stepFunc)
}

// items

func AddStep(ctx *godog.ScenarioContext, expr string, stepFunc any) {
	ctx.Step(fmt.Sprintf("^user adds %s with the following data$", expr), stepFunc)
}

func EditStep(ctx *godog.ScenarioContext, expr string, stepFunc any) {
	ctx.Step(fmt.Sprintf("^user edits %s with the following data$", expr), stepFunc)
}

func RemoveStep(ctx *godog.ScenarioContext, expr string, stepFunc any) {
	ctx.Step(fmt.Sprintf("^user removes %s with the following data$", expr), stepFunc)
}

// generic steps

func DefineStep(s *SharedData, ctx *godog.ScenarioContext, expr string) {
	ctx.Step(fmt.Sprintf("^user defines %s with the following data$", expr), SyncToDataMapStep(s))
}

func ReturnedErrorStep(ctx *godog.ScenarioContext, s *SharedData) {
	ctx.Step("^user should get the following error$", AssertErrorStep(s))
}
