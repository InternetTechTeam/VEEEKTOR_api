package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"

	"VEEEKTOR_api/pkg/database/pgsql"
	e "VEEEKTOR_api/pkg/errors"
)

var (
	AccessTokenLifeTime  = time.Minute * 15    // 15 minutes
	RefreshTokenLifeTime = time.Minute * 43800 // 30 days
	AccessKey            = []byte(os.Getenv("JWT_KEY"))
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Errors: -
func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", b), nil
}

// Errors: -
func GenerateAccessToken(userId, roleId, groupId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
		"exp":      time.Now().Add(AccessTokenLifeTime).Unix(),
		"user_id":  userId,
		"role_id":  roleId,
		"group_id": groupId,
	})

	tokenString, err := token.SignedString(AccessKey)
	if err != nil {
		log.Fatal(err)
	}

	return tokenString, nil
}

type RefreshToken struct {
	Token string `json:"refresh_token"`
}

// Errors: ErrTokenNotProvided
func GetRefreshTokenFromCookieOrBody(r *http.Request) (string, error) {
	// Try to get token from cookie
	var err error
	var rt RefreshToken
	var cookie *http.Cookie
	if cookie, err = r.Cookie("refresh_token"); err != nil &&
		!errors.Is(err, http.ErrNoCookie) {
		log.Fatal(err)
	}
	if !errors.Is(err, http.ErrNoCookie) {
		rt.Token = cookie.Value
	}

	// If cookie is missing, try to get it from body
	if errors.Is(err, http.ErrNoCookie) {
		bytes := make([]byte, r.ContentLength)
		r.Body.Read(bytes)

		if err = json.Unmarshal(
			bytes, &rt); err != nil {
			return "", e.ErrTokenNotProvided
		}
	}
	if rt.Token == "" {
		return "", e.ErrTokenNotProvided
	}
	return rt.Token, nil
}

// Errors: ErrTokenNotProvided, ErrTokenNotValid
func GetAccessTokenFromHeader(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", e.ErrTokenNotProvided
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", e.ErrTokenNotValid
	}

	if len(headerParts[1]) == 0 {
		return "", e.ErrTokenNotValid
	}
	return headerParts[1], nil
}

// Also checks authorization
// Errors: ErrTokenNotValid, ErrTokenExpired
func GetTokenClaims(accessToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(accessToken,
		func(token *jwt.Token) (i interface{}, err error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, e.ErrTokenNotValid
			}
			return AccessKey, nil
		})
	if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorExpired != 0 {
			return nil, e.ErrTokenExpired
		}
	}
	if err != nil {
		return nil, e.ErrTokenNotValid
	}

	claims := token.Claims.(jwt.MapClaims)

	// Convert json number to golang int
	claims["user_id"] = int(claims["user_id"].(float64))
	claims["role_id"] = int(claims["role_id"].(float64))
	claims["group_id"] = int(claims["group_id"].(float64))

	return claims, nil
}

// Errors: -
func IsRefreshTokenExpired(refreshToken string) (bool, error) {
	stmt, err := pgsql.DB.Prepare(
		`SELECT expires_at FROM sessions WHERE refresh_token=$1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	var expiresAt int64
	if err := stmt.QueryRow(&refreshToken).Scan(&expiresAt); err != nil {
		log.Fatal(err)
	}

	if expiresAt > time.Now().Unix() {
		return false, nil
	}

	return true, nil
}
