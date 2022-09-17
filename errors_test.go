package aeolic

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIError_Unwrap(t *testing.T) {
	appErr := APIError{
		StatusCode: http.StatusBadRequest,
		StatusText: http.StatusText(http.StatusBadRequest),
		Message:    "fla",
	}
	want := fmt.Errorf("%w", &appErr)
	got := appErr.Unwrap()

	assert.Error(t, want, got)
}
