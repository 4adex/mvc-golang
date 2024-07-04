package middleware

import (
	"context"
	// "encoding/json"
	// "fmt"
	"net/http"

	"github.com/4adex/mvc-golang/pkg/jwtutils"
	"github.com/4adex/mvc-golang/pkg/messages"
)


// func jsonResponse(w http.ResponseWriter, status int, redirect string) {
// 	response := map[string]string{
// 		"redirect": redirect,
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(status)
// 	json.NewEncoder(w).Encode(response)
// }

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
                // jsonResponse(w, http.StatusInternalServerError, "/viewbooks")
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
            messages.SetFlash(w, r, "Unauthorized Access", "error")
            http.Redirect(w, r, "/signin", http.StatusSeeOther)
            return
        }

        ctx := context.WithValue(r.Context(), "username", claims.Username)
        ctx = context.WithValue(ctx, "email", claims.Email)
        ctx = context.WithValue(ctx, "role", claims.Role)
        ctx = context.WithValue(ctx,"id",claims.Id)
        r = r.WithContext(ctx)

        next.ServeHTTP(w, r)
    })
}





func AdminMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        role := r.Context().Value("role").(string)
        // fmt.Println("Role is", role)
        if role != "admin" {
            messages.SetFlash(w, r, "Unauthorized Access", "error")
            http.Redirect(w, r, "/", http.StatusSeeOther)
            return
        }

        next.ServeHTTP(w, r)
    })
}




