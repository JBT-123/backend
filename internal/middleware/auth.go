package middleware

import (
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		// token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
		// 	//TODO: Set up JWT secret from config
		// 	return []byte("your_jwt_secret"), nil
		// })

		// if err != nil {
		// 	http.Error(w, "Invalid token", http.StatusUnauthorized)
		// 	return
		// }

		// Add user info to request context
		// claims := token.Claims.(jwt.MapClaims)
		// TODO: Add user to context
		// r = r.WithContext(context.WithValue(r.Context(), "user", claims))

		next.ServeHTTP(w, r)
	})
}
