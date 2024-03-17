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
			return "", err
		}
	}
	if rt.Token == "" {
		return "", e.ErrTokenNotProvided
	}
	return rt.Token, nil
}

func IsAccessTokenExpired(accessToken string) (bool, error) {

	return false, nil
}

func IsRefreshTokenExpired(refreshToken string) (bool, error) {

	return false, nil
}
