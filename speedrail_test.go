package speedrail_test

import (
	"context"
	"github.com/Kansuler/speedrail"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type SpeedrailTestSuite struct {
	suite.Suite
}

func (suite *SpeedrailTestSuite) TestPlanExecute() {
	plan := speedrail.Plan[any, any]()
	_, _, err := plan.Execute(context.Background(), nil, nil)
	suite.Error(err)
	suite.Equal(err.StatusCode(), http.StatusInternalServerError)
	suite.ErrorIs(err, speedrail.ErrNoStrategy)

	plan = speedrail.Plan[any, any](
		func(ctx context.Context, container any, model any) (context.Context, any, speedrail.Error) {
			return nil, nil, nil
		},
	)
	_, _, err = plan.Execute(context.Background(), nil, nil)
	suite.Error(err)
	suite.Equal(err.StatusCode(), http.StatusInternalServerError)
	suite.ErrorIs(err, speedrail.ErrNoContextReturned)
}

func TestSpeedrailTestSuite(t *testing.T) {
	suite.Run(t, new(SpeedrailTestSuite))
}
