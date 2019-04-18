package jwtutil

import (
    "log"
    "time"

    "github.com/dgrijalva/jwt-go"
    "github.com/damp_donkeys/internal/pkg/confidante"
)

// initially empty string populated on first token creation
var signingKey = ""

type CustomClaims struct {
    User string `json:"user"`
    jwt.StandardClaims
}

func populateSigningKey() error {
    if signingKey == "" {
        log.Println("Setting JWT signing key...")

        key, err := confidante.JWTSingingKey()

        if err != nil { return err }

        signingKey = key
        log.Println("JWT signing Key set!")
    } else {
        log.Println("JWT singing key already set, doing nothing.")
    }
    return nil
}

func CreateToken(user string, minutes int) (string, error) {
    err := populateSigningKey()

    if err != nil {
        return "", err
    }

    // Create the Claims
    claims := CustomClaims{
        user,
        jwt.StandardClaims{
            ExpiresAt: time.Now().Add(time.Minute * time.Duration(minutes)).Unix(),
            Issuer:    "go-csrcint",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    ss, err := token.SignedString([]byte(signingKey))
    return ss, err
}

func RefreshToken(tokenString string, minutes int) (string, error) {
    claims, err := ParseClaims(tokenString)
    if err != nil {
        return "", err
    }

    return CreateToken(claims.User, minutes)
}

func ParseClaims(tokenString string) (*CustomClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(signingKey), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
        return claims, nil
    } else {
        return nil, err
    }

}

func IsValidToken(givenToken string) (bool, error){
    token, err := jwt.Parse(givenToken, func(token *jwt.Token) (interface{}, error) {
        return []byte(signingKey), nil
    })

    return err == nil && token.Valid, err
}
