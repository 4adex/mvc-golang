package middleware

import (
	"context"
	"net/http"
	"github.com/4adex/mvc-golang/pkg/jwtutils"
    "strings"
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

        if strings.HasPrefix(r.URL.Path,"/admin") {
            if claims.Role != "admin" {
                http.Redirect(w, r, "/", http.StatusSeeOther)
                return
            }

        }

        ctx := context.WithValue(r.Context(), "username", claims.Username)
        ctx = context.WithValue(ctx, "email", claims.Email)
        ctx = context.WithValue(ctx, "role", claims.Role)
        ctx = context.WithValue(ctx,"id",claims.Id)
        r = r.WithContext(ctx)

        next.ServeHTTP(w, r)
    })
}





// func AdminMiddleware(next http.Handler) http.Handler {
//     return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//         role := r.Context().Value("role").(string)
//         if role != "admin" {
//             http.Redirect(w, r, "/", http.StatusSeeOther)
//             return
//         }

//         next.ServeHTTP(w, r)
//     })
// }




