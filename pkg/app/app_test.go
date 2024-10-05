package app_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"kycelis.dev/core/pkg/app"
	"kycelis.dev/core/pkg/response"
)

func TestNew(t *testing.T) {
	// Arrange.
	t.Parallel()

	logger := zap.NewNop()

	server := httptest.NewServer(app.New(logger).Server.Handler)
	defer server.Close()

	expected := response.NewStatusFromError(response.ErrServiceHealthy)

	// Act.
	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, server.URL+"/health", nil)
	require.NoError(t, err, "should not fail to create a new request")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err, "should not fail to send a request")

	defer res.Body.Close()

	// Assert.
	assert.Equal(t, expected.Code, res.StatusCode, "should return a 200 status code")

	status := new(response.Status)

	err = json.NewDecoder(res.Body).Decode(status)
	require.NoError(t, err, "should not fail to decode the response body")

	assert.Equal(t, *expected, *status, "should return the expected status")
}

func TestNewUnknownEndpoint(t *testing.T) {
	// Arrange.
	t.Parallel()

	logger := zap.NewNop()

	server := httptest.NewServer(app.New(logger).Server.Handler)
	defer server.Close()

	expected := response.NewStatusFromError(response.ErrUnknownEndpoint)

	// Act.
	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, server.URL+"/unknown", nil)
	require.NoError(t, err, "should not fail to create a new request")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err, "should not fail to send a request")

	defer res.Body.Close()

	// Assert.
	assert.Equal(t, expected.Code, res.StatusCode, "should return a 404 status code")

	status := new(response.Status)

	err = json.NewDecoder(res.Body).Decode(status)
	require.NoError(t, err, "should not fail to decode the response body")

	assert.Equal(t, *expected, *status, "should return the expected status")
}

func TestGetPort(t *testing.T) {
	// Arrange.
	logger := zap.NewNop()

	cases := []struct {
		name     string
		env      string
		expected int64
	}{
		{
			name:     "ValidPort",
			env:      "9000",
			expected: 9000,
		},
		{
			name:     "InvalidPort",
			env:      "invalid",
			expected: app.DefaultPort,
		},
		{
			name:     "EmptyPort",
			env:      "",
			expected: app.DefaultPort,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(run *testing.T) {
			// Arrange.
			run.Setenv("PORT", testCase.env)

			// Act.
			port := app.GetPort(logger)

			// Assert.
			require.Equal(run, testCase.expected, port, "should return the expected port")
		})
	}
}
