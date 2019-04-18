package router

import(
	"log"
	"net/http"
	"math/rand"
	"crypto/sha1"
	"encoding/base64"

	"github.com/damp_donkeys/internal/pkg/jwtutil"
)

const(
	JWTDuration = 30 // In minutes. How long the token will stay valid between requests
	UserCodeLength = 5 // How long user / company unique codes should be
	UserCodeChars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // The characters allowed in the unique codes
	DBName = "dev" // Which database to populate during API requests
)

func hashHelper(str string) string {
	bv := []byte(str)
	hasher := sha1.New()
	hasher.Write(bv)
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func genCodeHelper(codeLength int) string {
	// 26 lowercase + 26 uppercase + 10 digit possibilities
	bv := make([]byte, codeLength)

	for i := range bv {
		bv[i] = UserCodeChars[rand.Int63() % int64(len(UserCodeChars))]
	}

	return string(bv)
}

// Given a JWT and user, verify that the JWT is valid, and belongs to that user
// returns
//          appropriate http status code, nil
//          0, refreshed JWT
func tokenRefreshHelper(jwt string, user string, minutes int) (int, string) {
    // -> JWT ERROR HANDLING
    old_jwt := jwt
    is_valid, err := jwtutil.IsValidToken(old_jwt)
    if !is_valid {
        log.Printf("IsValidToken error: %s\n", err)
        return http.StatusUnauthorized, ""
    }

    claims, err := jwtutil.ParseClaims(old_jwt)

    // Something went wrong internally
    if err != nil {
        log.Printf("ParseClaims error: %s\n", err)
        return http.StatusInternalServerError, ""
    }

    jwt_user := claims.User
    // Only admins should be able to add a company
    if jwt_user != user {
        // Note: Hard coded "admin" could (/should) eventually be replaced with a cross check to some 'Admins' Table in the db 
        log.Printf("JWT invalid for requested user [%s != %s]\n", jwt_user, user)
        return http.StatusUnauthorized, ""
    }

    new_jwt, err := jwtutil.RefreshToken(old_jwt, minutes)

    if err != nil {
        log.Printf("RefreshToken error: \n", err)
        return http.StatusInternalServerError, ""
    }
    // <- END JWT ERROR HANDLING

    return 0, new_jwt
}
