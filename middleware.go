package ermeslog

import (
    "context"
    "net/http"
    "time"
    "github.com/google/uuid"
)

type contextKey string

const (
    RequestIDKey contextKey = "requestID"
    UserIDKey    contextKey = "userID"
    TraceIDKey   contextKey = "traceID"
)

// HTTPMiddleware adds logging to HTTP requests
func (l *Logger) HTTPMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Generate request ID if not present
        requestID := r.Header.Get("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }
        
        // Add to context
        ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
        r = r.WithContext(ctx)
        
        // Log request
        l.WithContext(ctx).WithFields(map[string]interface{}{
            "method": r.Method,
            "path":   r.URL.Path,
            "ip":     r.RemoteAddr,
        }).Info("Request started")
        
        // Wrap ResponseWriter to capture status
        wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
        
        next.ServeHTTP(wrapped, r)
        
        // Log response
        l.WithContext(ctx).WithFields(map[string]interface{}{
            "method":   r.Method,
            "path":     r.URL.Path,
            "status":   wrapped.statusCode,
            "duration": time.Since(start).Milliseconds(),
        }).Info("Request completed")
    })
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}