package speedrail_test

import (
	"encoding/json"
	"errors"
	"github.com/Kansuler/speedrail"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type SpeedrailErrorTestSuite struct {
	suite.Suite
}

var error1 = errors.New("error 1")

func errorTestFunction1() speedrail.Error {
	return speedrail.NewError(error1, http.StatusInternalServerError, "error 1 message")
}

var error2 = errors.New("error 2")

func errorTestFunction2() speedrail.Error {
	return speedrail.NewError(error2, http.StatusInternalServerError, "error 2 message")
}

var error3 = errors.New("error 3")

func (suite *SpeedrailErrorTestSuite) TestMerge() {
	err1 := errorTestFunction1()
	err2 := errorTestFunction2()
	err3 := errorTestFunction2()

	result := err1.Merge(err2).Merge(err3)
	suite.Equal(3, len(result.Trail()))
	suite.Equal(http.StatusInternalServerError, result.StatusCode())
	suite.Equal("error 1 message; error 2 message; error 2 message", result.Error())
	suite.True(errors.Is(result, error1))
	suite.True(errors.Is(result, error2))
	suite.False(errors.Is(result, error3))
	b, err := json.Marshal(result.Trail())
	suite.NoError(err)
	suite.JSONEq(`{"[1]speedrail_test.errorTestFunction1":"error 1","[2]speedrail_test.errorTestFunction2":"error 2","[3]speedrail_test.errorTestFunction2":"error 2"}`, string(b))

	err4 := speedrail.NewError(nil, http.StatusInternalServerError, "error 3 message")
	b, err = json.Marshal(err4.Trail())
	suite.NoError(err)
	suite.JSONEq(`{"[1]speedrail_test.(*SpeedrailErrorTestSuite).TestMerge":"error 3 message"}`, string(b))
}

func TestSpeedrailErrorTestSuite(t *testing.T) {
	suite.Run(t, new(SpeedrailErrorTestSuite))
}
