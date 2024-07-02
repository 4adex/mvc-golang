package middleware

import (
    "context"
    "net/http"
    "github.com/4adex/mvc-golang/pkg/jwtutils"
)


func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Skip authentication for sign-in and sign-up pages
        if r.URL.Path == "/signin" || r.URL.Path == "/signup" {
            next.ServeHTTP(w, r)
            return
        }

        cookie, err := r.Cookie("token")
        if err != nil {
            if err == http.ErrNoCookie {
                // Redirect to sign-in page
                http.Redirect(w, r, "/signin", http.StatusSeeOther)
                return
            }
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        tokenStr := cookie.Value
        claims, err := jwtutils.ValidateJWT(tokenStr)
        if err != nil {
            // Redirect to sign-in page
            http.Redirect(w, r, "/signin", http.StatusSeeOther)
            return
        }

        ctx := context.WithValue(r.Context(), "username", claims.Username)
        ctx = context.WithValue(ctx, "email", claims.Email)
        ctx = context.WithValue(ctx, "role", claims.Role)
        r = r.WithContext(ctx)

        next.ServeHTTP(w, r)
    })
}





func AdminMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tokenCookie, err := r.Cookie("token")
        if err != nil {
            if err == http.ErrNoCookie {
                w.WriteHeader(http.StatusUnauthorized)
                return
            }
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        tokenStr := tokenCookie.Value
        claims, err := jwtutils.ValidateJWT(tokenStr)
        if err != nil {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }

        if claims.Role != "admin" {
            w.WriteHeader(http.StatusForbidden)
            return
        }

        next.ServeHTTP(w, r)
    })
}

