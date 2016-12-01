package negronichaos

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// Assert we get an expected number of "bad" responses given a known random
// seed and desired frequency
func TestMiddleware(t *testing.T) {
	mw := NewMiddleware(1234, 0.64)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Good")
	}

	var goodResponses int
	var badResponses int

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)

	for i := 0; i < 100; i++ {
		w := httptest.NewRecorder()
		mw.ServeHTTP(w, req, dummyHandler)
		if w.Code == http.StatusOK {
			goodResponses++
		} else {
			badResponses++
		}
	}
	require.Equal(t, 34, goodResponses)
	require.Equal(t, badResponses, 66)

}
