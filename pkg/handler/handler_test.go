package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteCode(t *testing.T) {
	// Form a dummy request with RAW data
	req, err := http.NewRequest("POST", "/execute", strings.NewReader("print \"Hello, World!\";"))
	assert.NoError(t, err, "Creating request should not produce an error")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ExecuteCode)

	// Call the handler function
	handler.ServeHTTP(rr, req)

	// Check the status code and response body using testify assertions
	assert.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code")
	assert.Equal(t, "Output: Hello, World!\n", rr.Body.String(), "Handler returned unexpected body")
}
