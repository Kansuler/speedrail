package speedrail_test

import (
	"context"
	"errors"
	"github.com/Kansuler/speedrail"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type SpeedrailStrategyTestSuite struct {
	suite.Suite
}

type strategyTestModel struct {
	CriteriaMet bool
}

func (suite *SpeedrailStrategyTestSuite) TestIf() {
	plan := speedrail.Plan(
		speedrail.If(
			func(model strategyTestModel) bool {
				return true
			},
			func(ctx context.Context, container any, model strategyTestModel) (context.Context, strategyTestModel, speedrail.Error) {
				model.CriteriaMet = true
				return ctx, model, nil
			}),
	)

	_, model, err := plan.Execute(context.Background(), nil, strategyTestModel{})
	suite.NoError(err)
	suite.True(model.CriteriaMet)

	plan = speedrail.Plan(
		speedrail.If(
			func(model strategyTestModel) bool {
				return false
			},
			func(ctx context.Context, container any, model strategyTestModel) (context.Context, strategyTestModel, speedrail.Error) {
				model.CriteriaMet = false
				return ctx, model, nil
			}),
	)

	_, model, err = plan.Execute(context.Background(), nil, strategyTestModel{CriteriaMet: true})
	suite.NoError(err)
	suite.True(model.CriteriaMet)
}

func (suite *SpeedrailStrategyTestSuite) TestIfElse() {
	plan := speedrail.Plan(
		speedrail.IfElse(
			func(model strategyTestModel) bool {
				return true
			},
			func(ctx context.Context, container any, model strategyTestModel) (context.Context, strategyTestModel, speedrail.Error) {
				model.CriteriaMet = true
				return ctx, model, nil
			},
			func(ctx context.Context, container any, model strategyTestModel) (context.Context, strategyTestModel, speedrail.Error) {
				model.CriteriaMet = false
				return ctx, model, nil
			},
		),
	)

	_, model, err := plan.Execute(context.Background(), nil, strategyTestModel{})
	suite.NoError(err)
	suite.True(model.CriteriaMet)

	plan = speedrail.Plan(
		speedrail.IfElse(
			func(model strategyTestModel) bool {
				return false
			},
			func(ctx context.Context, container any, model strategyTestModel) (context.Context, strategyTestModel, speedrail.Error) {
				model.CriteriaMet = false
				return ctx, model, nil
			},
			func(ctx context.Context, container any, model strategyTestModel) (context.Context, strategyTestModel, speedrail.Error) {
				model.CriteriaMet = true
				return ctx, model, nil
			},
		),
	)

	_, model, err = plan.Execute(context.Background(), nil, strategyTestModel{CriteriaMet: true})
	suite.NoError(err)
	suite.True(model.CriteriaMet)
}

func (suite *SpeedrailStrategyTestSuite) TestMerge() {
	plan := speedrail.Plan(
		speedrail.Merge(
			func(ctx context.Context, container any, model strategyTestModel) (context.Context, strategyTestModel, speedrail.Error) {
				return ctx, model, speedrail.NewError(errors.New("error 1"), http.StatusBadRequest, "error 1")
			},
			func(ctx context.Context, container any, model strategyTestModel) (context.Context, strategyTestModel, speedrail.Error) {
				return ctx, model, speedrail.NewError(errors.New("error 2"), http.StatusForbidden, "error 2")
			},
			func(ctx context.Context, container any, model strategyTestModel) (context.Context, strategyTestModel, speedrail.Error) {
				model.CriteriaMet = true
				return ctx, model, nil
			},
		),
	)

	_, model, err := plan.Execute(context.Background(), nil, strategyTestModel{})
	suite.Error(err)
	suite.Equal(http.StatusForbidden, err.StatusCode())
	suite.Equal("error 1; error 2", err.Error())
	suite.True(model.CriteriaMet)
}

func (suite *SpeedrailStrategyTestSuite) TestGroup() {
	plan := speedrail.Plan(
		speedrail.Group(
			func(ctx context.Context, container any, model strategyTestModel) (context.Context, strategyTestModel, speedrail.Error) {
				return ctx, model, speedrail.NewError(errors.New("error 1"), http.StatusBadRequest, "error 1")
			},
			func(ctx context.Context, container any, model strategyTestModel) (context.Context, strategyTestModel, speedrail.Error) {
				model.CriteriaMet = false
				return ctx, model, nil
			},
		),
	)

	_, model, err := plan.Execute(context.Background(), nil, strategyTestModel{CriteriaMet: true})
	suite.Error(err)
	suite.Equal(http.StatusBadRequest, err.StatusCode())
	suite.Equal("error 1", err.Error())
	suite.True(model.CriteriaMet)
}

func (suite *SpeedrailStrategyTestSuite) TestThrowError() {
	plan := speedrail.Plan[any, strategyTestModel](
		speedrail.ThrowError[any, strategyTestModel](speedrail.NewError(errors.New("error 1"), http.StatusBadRequest, "error 1")),
	)

	_, model, err := plan.Execute(context.Background(), nil, strategyTestModel{CriteriaMet: true})
	suite.Error(err)
	suite.Equal(http.StatusBadRequest, err.StatusCode())
	suite.Equal("error 1", err.Error())
	suite.True(model.CriteriaMet)
}

func TestSpeedrailStrategyTestSuite(t *testing.T) {
	suite.Run(t, new(SpeedrailStrategyTestSuite))
}
