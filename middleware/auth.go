package middleware

import (
    "fmt"
    "net/http"
    "strings"
    "context"
    "github.com/golang-jwt/jwt/v5"
)

// Define a key type to use when storing values in context
type contextKey string

const UserContextKey = contextKey("userID")

// Secret key used for signing the JWTs (you should store this securely!)
var jwtSecret = []byte("your-secret-key")

// AuthMiddleware verifies the JWT in the Authorization header
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        // 1. Get the Authorization header (e.g., "Bearer <token>")
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header missing", http.StatusUnauthorized)
            return
        }

        // 2. Check if the header starts with "Bearer "
        if !strings.HasPrefix(authHeader, "Bearer ") {
            http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
            return
        }

        // 3. Extract the token string by removing "Bearer " prefix
        tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

        // 4. Parse and validate the token
        token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
            // Ensure the signing method is what we expect
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method")
            }
            return jwtSecret, nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
            return
        }

        // 5. Extract user ID (or other claims) if needed
        if claims, ok := token.Claims.(jwt.MapClaims); ok {
            userID := claims["user_id"]

            // 6. Save the user ID in the request context so handlers can use it
            ctx := context.WithValue(r.Context(), UserContextKey, userID)

            // 7. Call the next handler with the new context
            next.ServeHTTP(w, r.WithContext(ctx))
        } else {
            http.Error(w, "Invalid token claims", http.StatusUnauthorized)
            return
        }
    })
}