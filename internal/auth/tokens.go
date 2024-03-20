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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
		"exp":     time.Now().Add(AccessTokenLifeTime).Unix(),
		"user_id": user_id,
		"role_id": role_id,
	})

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
	if cookie, err = r.Cookie("refresh_token"); err != nil &&
		!errors.Is(err, http.ErrNoCookie) {
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

func GetTokenClaims(accessToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(accessToken,
		func(token *jwt.Token) (i interface{}, err error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, e.ErrTokenNotValid
			}
			return AccessKey, nil
		})
	if err != nil {
		return nil, e.ErrTokenNotValid
	}

	claims := token.Claims.(jwt.MapClaims)

	return claims, nil
}

func IsAccessTokenExpired(accessToken string) (bool, error) {
	var claims jwt.MapClaims
	var err error
	if claims, err = GetTokenClaims(accessToken); err != nil {
		return true, err
	}

	// Cast json number to golang int
	exp := claims["exp"].(float64)
	if int64(exp) > time.Now().Unix() {
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

func GetUserIdFromRequest(r *http.Request) (int, error) {
	accessToken, err := GetAccessTokenFromHeader(r)
	if err != nil {
		return 0, err
	}

	claims, err := GetTokenClaims(accessToken)
	if err != nil {
		return 0, err
	}

	// Cast json number to golang int
	userId := int(claims["user_id"].(float64))

	return userId, nil
}

func CheckUserAuthorized(r *http.Request) (bool, error) {
	accessToken, err := GetAccessTokenFromHeader(r)
	if err != nil {
		return false, err
	}

	exp, err := IsAccessTokenExpired(accessToken)
	if err != nil {
		return false, err
	}

	if exp {
		return false, nil
	}

	return true, nil
}
