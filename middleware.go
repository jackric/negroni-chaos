package negronichaos

import (
	"math/rand"
	"net/http"
)

type Middleware struct {
	frequency float32
	rand      *rand.Rand
}

// NewMiddleware creates a chaotic middleware handler that will randomly
// respond with an Internal Server error with the given frequency (0.0 to 1.0)
func NewMiddleware(seed int64, frequency float32) *Middleware {
	return &Middleware{
		rand:      rand.New(rand.NewSource(seed)),
		frequency: frequency,
	}
}

func (l *Middleware) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	cutoff := 1.0 - l.frequency
	sample := l.rand.Float32()
	if sample > cutoff {
		http.Error(w, "Why so serious?", http.StatusInternalServerError)
	} else {
		// Call the next middleware handler
		next(w, req)
	}
}
