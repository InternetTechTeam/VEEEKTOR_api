package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"

	"VEEEKTOR_api/pkg/database/pgsql"
	e "VEEEKTOR_api/pkg/errors"
)

var (
	AccessTokenLifeTime  = time.Minute * 60    // One hour
	RefreshTokenLifeTime = time.Minute * 43800 // 30 days
	// Should be implemented via environment
	AccessKey = []byte(os.Getenv("JWT_KEY"))
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

func GenerateAccessToken(user_id int, role_id int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(AccessTokenLifeTime).Unix()
	claims["user_id"] = user_id
	claims["role_id"] = role_id
	tokenString, err := token.SignedString(AccessKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

type RefreshToken struct {
	Token string `json:"refresh_token"`
}

func GetRefreshTokenFromCookieOrBody(r *http.Request) (string, error) {
	// Try to get token from cookie
	var err error
	var rt RefreshToken
	var cookie *http.Cookie
	if cookie, err = r.Cookie(
		"refresh_token"); err != nil && !errors.Is(err, http.ErrNoCookie) {
		log.Print(err)
		return "", err
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

func GetTokenClaims(accessToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, e.ErrTokenNotValid
		}
		return AccessKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, e.ErrTokenNotValid
	}

	return claims, nil
}

func IsAccessTokenExpired(accessToken string) (bool, error) {
	var claims jwt.MapClaims
	var err error
	if claims, err = GetTokenClaims(accessToken); err != nil {
		return true, err
	}

	if claims["exp"].(int64) > time.Now().Unix() {
		return false, nil
	}

	return true, nil
}

func IsRefreshTokenExpired(refreshToken string) (bool, error) {
	stmt, err := pgsql.DB.Prepare(
		`SELECT expires_at FROM sessions WHERE refresh_token=$1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	var expiresAt int64
	if err := stmt.QueryRow(&refreshToken).Scan(&expiresAt); err != nil {
		log.Print(err)
		return true, err
	}

	if expiresAt > time.Now().Unix() {
		return false, nil
	}

	return true, nil
}
