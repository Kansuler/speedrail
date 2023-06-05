package speedrail_test

import (
	"context"
	"github.com/Kansuler/speedrail"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SpeedrailConditionTestSuite struct {
	suite.Suite
}

type conditionTestModel struct {
	CriteriaMet bool
}

func (suite *SpeedrailConditionTestSuite) TestAnd() {
	plan := speedrail.Plan(
		speedrail.If(
			speedrail.And(
				func(model conditionTestModel) bool {
					return true
				},
				func(model conditionTestModel) bool {
					return true
				},
			),
			func(ctx context.Context, container any, model conditionTestModel) (context.Context, conditionTestModel, speedrail.Error) {
				model.CriteriaMet = true
				return ctx, model, nil
			},
		),
	)

	_, model, err := plan.Execute(context.Background(), nil, conditionTestModel{})
	suite.NoError(err)
	suite.True(model.CriteriaMet)

	plan = speedrail.Plan(
		speedrail.If(
			speedrail.And(
				func(model conditionTestModel) bool {
					return false
				},
				func(model conditionTestModel) bool {
					return true
				},
			),
			func(ctx context.Context, container any, model conditionTestModel) (context.Context, conditionTestModel, speedrail.Error) {
				model.CriteriaMet = false
				return ctx, model, nil
			},
		),
	)

	_, model, err = plan.Execute(context.Background(), nil, conditionTestModel{CriteriaMet: true})
	suite.NoError(err)
	suite.True(model.CriteriaMet)

	plan = speedrail.Plan(
		speedrail.If(
			speedrail.And(
				func(model conditionTestModel) bool {
					return true
				},
				func(model conditionTestModel) bool {
					return false
				},
			),
			func(ctx context.Context, container any, model conditionTestModel) (context.Context, conditionTestModel, speedrail.Error) {
				model.CriteriaMet = false
				return ctx, model, nil
			},
		),
	)

	_, model, err = plan.Execute(context.Background(), nil, conditionTestModel{CriteriaMet: true})
	suite.NoError(err)
	suite.True(model.CriteriaMet)
}

func (suite *SpeedrailConditionTestSuite) TestOr() {
	plan := speedrail.Plan(
		speedrail.If(
			speedrail.Or(
				func(model conditionTestModel) bool {
					return true
				},
				func(model conditionTestModel) bool {
					return true
				},
			),
			func(ctx context.Context, container any, model conditionTestModel) (context.Context, conditionTestModel, speedrail.Error) {
				model.CriteriaMet = true
				return ctx, model, nil
			},
		),
	)

	_, model, err := plan.Execute(context.Background(), nil, conditionTestModel{})
	suite.NoError(err)
	suite.True(model.CriteriaMet)

	plan = speedrail.Plan(
		speedrail.If(
			speedrail.Or(
				func(model conditionTestModel) bool {
					return false
				},
				func(model conditionTestModel) bool {
					return true
				},
			),
			func(ctx context.Context, container any, model conditionTestModel) (context.Context, conditionTestModel, speedrail.Error) {
				model.CriteriaMet = true
				return ctx, model, nil
			},
		),
	)

	_, model, err = plan.Execute(context.Background(), nil, conditionTestModel{})
	suite.NoError(err)
	suite.True(model.CriteriaMet)

	plan = speedrail.Plan(
		speedrail.If(
			speedrail.Or(
				func(model conditionTestModel) bool {
					return false
				},
				func(model conditionTestModel) bool {
					return true
				},
			),
			func(ctx context.Context, container any, model conditionTestModel) (context.Context, conditionTestModel, speedrail.Error) {
				model.CriteriaMet = true
				return ctx, model, nil
			},
		),
	)

	_, model, err = plan.Execute(context.Background(), nil, conditionTestModel{})
	suite.NoError(err)
	suite.True(model.CriteriaMet)

	plan = speedrail.Plan(
		speedrail.If(
			speedrail.Or(
				func(model conditionTestModel) bool {
					return false
				},
				func(model conditionTestModel) bool {
					return false
				},
			),
			func(ctx context.Context, container any, model conditionTestModel) (context.Context, conditionTestModel, speedrail.Error) {
				model.CriteriaMet = false
				return ctx, model, nil
			},
		),
	)

	_, model, err = plan.Execute(context.Background(), nil, conditionTestModel{CriteriaMet: true})
	suite.NoError(err)
	suite.True(model.CriteriaMet)
}

func (suite *SpeedrailConditionTestSuite) TestNot() {
	plan := speedrail.Plan(
		speedrail.If(
			speedrail.And(
				speedrail.Not(func(model conditionTestModel) bool {
					return false
				}),
				func(model conditionTestModel) bool {
					return true
				},
			),
			func(ctx context.Context, container any, model conditionTestModel) (context.Context, conditionTestModel, speedrail.Error) {
				model.CriteriaMet = true
				return ctx, model, nil
			},
		),
	)

	_, model, err := plan.Execute(context.Background(), nil, conditionTestModel{})
	suite.NoError(err)
	suite.True(model.CriteriaMet)
}

func (suite *SpeedrailConditionTestSuite) TestMixed() {
	plan := speedrail.Plan(
		speedrail.If(
			speedrail.And(
				speedrail.Not(func(model conditionTestModel) bool {
					return false
				}),
				speedrail.Or(
					speedrail.Not(func(model conditionTestModel) bool {
						return false
					}),
					func(model conditionTestModel) bool {
						return false
					},
				),
			),
			func(ctx context.Context, container any, model conditionTestModel) (context.Context, conditionTestModel, speedrail.Error) {
				model.CriteriaMet = true
				return ctx, model, nil
			},
		),
	)

	_, model, err := plan.Execute(context.Background(), nil, conditionTestModel{})
	suite.NoError(err)
	suite.True(model.CriteriaMet)
}

func TestSpeedrailConditionTestSuite(t *testing.T) {
	suite.Run(t, new(SpeedrailConditionTestSuite))
}
