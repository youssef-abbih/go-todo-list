package middleware

import "net/http"

// secureHeadersMiddleware is an HTTP middleware that adds basic security headers to each response.
// 
// These headers help protect the application against common web vulnerabilities:
//
// - X-XSS-Protection: Enables the browser's built-in cross-site scripting (XSS) filter.
// - X-Frame-Options: Prevents the page from being embedded in an iframe (clickjacking defense).
// - Content-Security-Policy: Restricts the sources of content (default to 'self') to mitigate injection attacks.
// - Referrer-Policy: Controls the information sent in the Referer header (set to 'no-referrer').
//
// Use this middleware early in the middleware chain to ensure headers are always applied.
func SecureHeadersMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("X-XSS-Protection", "1; mode-block")
    w.Header().Set("X-Frame-Options", "deny")
	w.Header().Set("Content-Security-Policy", "default-src 'self'")
	w.Header().Set("Referrer-Policy", "no-referrer")

    next.ServeHTTP(w, r)
  })
}