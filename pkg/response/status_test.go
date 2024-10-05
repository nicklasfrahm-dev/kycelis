package response_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"kycelis.dev/core/pkg/response"
)

var ErrTestError = errors.New("TestError")

func TestNewStatus(t *testing.T) {
	// Arrange.
	t.Parallel()

	expected := response.Status{
		Code:    http.StatusOK,
		Title:   http.StatusText(http.StatusOK),
		Message: "Example",
	}

	// Act.
	status := response.NewStatus(expected.Code, expected.Message)

	// Assert.
	assert.Equal(t, expected, *status, "should return the correct status")
}

func TestNewStatusFromError(t *testing.T) {
	// Arrange.
	t.Parallel()

	cases := []struct {
		err      error
		expected response.Status
	}{
		{
			err:      ErrTestError,
			expected: *response.NewStatus(http.StatusInternalServerError, response.ErrUnexpectedError.Error()),
		},
		{
			err:      response.ErrUnknownEndpoint,
			expected: *response.NewStatus(http.StatusNotFound, response.ErrUnknownEndpoint.Error()),
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.err.Error(), func(run *testing.T) {
			// Arrange.
			run.Parallel()

			// Act.
			status := response.NewStatusFromError(testCase.err)

			// Assert.
			assert.Equal(run, testCase.expected, *status, "should return the correct status")
		})
	}
}
