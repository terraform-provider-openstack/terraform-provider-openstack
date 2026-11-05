package openstack

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newRequestTimeoutTestServer returns a server that answers the first
// `failures` requests with the given status code and then 200.
func newRequestTimeoutTestServer(t *testing.T, failures int32, code int) (*httptest.Server, *atomic.Int32) {
	t.Helper()

	var calls atomic.Int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if calls.Add(1) <= failures {
			w.WriteHeader(code)
			_, _ = w.Write([]byte("<html><body><h1>error</h1></body></html>"))

			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(srv.Close)

	return srv, &calls
}

func newRequestTimeoutTestClient(maxRetries uint) *gophercloud.ProviderClient {
	return &gophercloud.ProviderClient{
		HTTPClient: http.Client{},
		RetryFunc:  requestTimeoutRetryFunc(maxRetries, time.Millisecond),
	}
}

// doTestRequest makes a request and closes the response body, returning only
// the error.
func doTestRequest(client *gophercloud.ProviderClient, method, url string) error {
	resp, err := client.Request(context.Background(), method, url, &gophercloud.RequestOpts{})
	if resp != nil {
		resp.Body.Close()
	}

	return err
}

func TestUnitRequestTimeoutRetryFuncRetriesGet(t *testing.T) {
	srv, calls := newRequestTimeoutTestServer(t, 2, http.StatusRequestTimeout)
	client := newRequestTimeoutTestClient(3)

	err := doTestRequest(client, http.MethodGet, srv.URL)
	require.NoError(t, err)
	assert.Equal(t, int32(3), calls.Load())
}

func TestUnitRequestTimeoutRetryFuncGivesUp(t *testing.T) {
	srv, calls := newRequestTimeoutTestServer(t, 10, http.StatusRequestTimeout)
	client := newRequestTimeoutTestClient(3)

	err := doTestRequest(client, http.MethodGet, srv.URL)
	require.Error(t, err)
	assert.True(t, gophercloud.ResponseCodeIs(err, http.StatusRequestTimeout))
	// One initial request plus three retries.
	assert.Equal(t, int32(4), calls.Load())
}

func TestUnitRequestTimeoutRetryFuncSkipsNonIdempotent(t *testing.T) {
	srv, calls := newRequestTimeoutTestServer(t, 1, http.StatusRequestTimeout)
	client := newRequestTimeoutTestClient(3)

	err := doTestRequest(client, http.MethodPost, srv.URL)
	require.Error(t, err)
	assert.True(t, gophercloud.ResponseCodeIs(err, http.StatusRequestTimeout))
	assert.Equal(t, int32(1), calls.Load())
}

func TestUnitRequestTimeoutRetryFuncSkipsOtherCodes(t *testing.T) {
	srv, calls := newRequestTimeoutTestServer(t, 1, http.StatusNotFound)
	client := newRequestTimeoutTestClient(3)

	err := doTestRequest(client, http.MethodGet, srv.URL)
	require.Error(t, err)
	assert.True(t, gophercloud.ResponseCodeIs(err, http.StatusNotFound))
	assert.Equal(t, int32(1), calls.Load())
}

func TestUnitRequestTimeoutRetryFuncContextCanceled(t *testing.T) {
	retryFunc := requestTimeoutRetryFunc(3, time.Hour)
	origErr := gophercloud.ErrUnexpectedResponseCode{Actual: http.StatusRequestTimeout}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := retryFunc(ctx, http.MethodGet, "http://example", nil, origErr, 1)
	require.Error(t, err)
	assert.True(t, gophercloud.ResponseCodeIs(err, http.StatusRequestTimeout))
}
