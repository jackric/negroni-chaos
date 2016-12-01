package negronichaos

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"time"

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

func TestSlowMiddleware(t *testing.T) {
	shortest := 20 * time.Millisecond
	longest := 80 * time.Millisecond
	// Allow for some time for the real handler underneath
	toleranceLongest := longest + (5 * time.Millisecond)
	mw := NewSlowMiddleware(1234, shortest, longest)
	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Good")
	}

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)

	for i := 0; i < 100; i++ {
		w := httptest.NewRecorder()
		begin := time.Now()
		mw.ServeHTTP(w, req, dummyHandler)
		end := time.Now()
		duration := end.Sub(begin)
		require.True(t, duration < toleranceLongest, "Reponse took too long, %s > %s", duration, toleranceLongest)
		require.True(t, duration > shortest, "Reponse was too fast, %s < %s", duration, shortest)
	}

}
