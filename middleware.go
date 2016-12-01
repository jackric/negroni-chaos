package negronichaos

import (
	"math/rand"
	"net/http"
	"time"
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

// SlowMiddleware randomly delays the response
type SlowMiddleware struct {
	shortest time.Duration
	longest  time.Duration
	rand     *rand.Rand
}

func NewSlowMiddleware(seed int64, shortest time.Duration, longest time.Duration) *SlowMiddleware {
	if longest < shortest {
		panic("Longest delay must be longer than shortest delay")
	}
	return &SlowMiddleware{
		rand:     rand.New(rand.NewSource(seed)),
		shortest: shortest,
		longest:  longest,
	}
}

func (mw *SlowMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	delta := mw.longest.Nanoseconds() - mw.shortest.Nanoseconds()
	delta = int64(float64(delta) * mw.rand.Float64())
	deltaDuration := time.Duration(delta) * time.Nanosecond
	time.Sleep(mw.shortest + deltaDuration)
	next(w, req)

}
