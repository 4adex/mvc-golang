package jwtutils

import (
    "os"
    "time"
    "github.com/dgrijalva/jwt-go"
    "github.com/joho/godotenv"
    "log"
)


func init() {
    if err := godotenv.Load(); err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }
}

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type Claims struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Role     string `json:"role"`
    Id string `json:"id"`
    jwt.StandardClaims
}

func GenerateJWT(username string, email string, role string, id string) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        Username: username,
        Email:    email,
        Role:     role,
        Id: id,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}

func ValidateJWT(tokenStr string) (*Claims, error) {
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })

    if err != nil {
        if err == jwt.ErrSignatureInvalid {
            return nil, err
        }
        return nil, err
    }

    if !token.Valid {
        return nil, err
    }

    return claims, nil
}
