package auth

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"VEEEKTOR_api/pkg/database/pgsql"
	e "VEEEKTOR_api/pkg/errors"
)

type Session struct {
	Id           int       `json:"id"`
	UserId       int       `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// User id, role_id and group_id must be valid.
// Errors: -
func StoreSession(user_id, role_id, groupId int) (TokenResponse, error) {
	stmt, err := pgsql.DB.Prepare(
		`INSERT INTO sessions (user_id, refresh_token, expires_at)
		VALUES ($1, $2, $3)`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	var resp TokenResponse
	resp.AccessToken, _ = GenerateAccessToken(user_id, role_id, groupId)
	resp.RefreshToken, _ = GenerateRefreshToken()

	CheckSessionsCount(user_id)

	if _, err := stmt.Exec(
		&user_id, &resp.RefreshToken,
		time.Now().Add(RefreshTokenLifeTime)); err != nil {
		log.Fatal(err)
	}

	return resp, nil
}

func UpdateSession(sess Session) (TokenResponse, error) {
	stmt, err := pgsql.DB.Prepare(
		`UPDATE sessions SET 
		refresh_token=$2, expires_at=$3
		WHERE refresh_token=$1`)
	if err != nil {
		log.Fatal(err)
	}

	var roleId, groupId int
	err = pgsql.DB.QueryRow(
		`SELECT role_id, group_id FROM users WHERE id=$1`,
		&sess.UserId).Scan(&roleId, &groupId)
	if err != nil {
		log.Fatal(err)
	}

	var resp TokenResponse
	resp.AccessToken, _ = GenerateAccessToken(sess.UserId, roleId, groupId)
	resp.RefreshToken, _ = GenerateRefreshToken()

	if _, err = stmt.Exec(
		&sess.RefreshToken, &resp.RefreshToken,
		time.Now().Add(RefreshTokenLifeTime)); err != nil {
		log.Fatal(err)
	}

	return resp, nil
}

// Errors: -
func DeleteSessionByRT(refreshToken string) error {
	stmt, err := pgsql.DB.Prepare(
		`DELETE FROM sessions WHERE refresh_token=$1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	if _, err := stmt.Exec(&refreshToken); err != nil {
		log.Fatal(err)
	}

	return nil
}

// Errors: -
func CheckSessionsCount(user_id int) {
	getStmt, err := pgsql.DB.Prepare(
		`SELECT COUNT(*) FROM sessions WHERE user_id=$1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}
	var count int
	if getStmt.QueryRow(&user_id).Scan(&count); count < 5 {
		return
	}

	deleteStmt, err := pgsql.DB.Prepare(
		`DELETE FROM sessions WHERE user_id = $1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	if _, err := deleteStmt.Exec(&user_id); err != nil {
		log.Fatal(err)
	}
}

// Errors: ErrSessionNotExist
func GetSessionByRefreshToken(refreshToken string) (Session, error) {
	stmt, err := pgsql.DB.Prepare(
		`SELECT id, user_id, expires_at 
		FROM sessions WHERE refresh_token=$1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	var sess Session
	if err := stmt.QueryRow(&refreshToken).Scan(
		&sess.Id, &sess.UserId, &sess.ExpiresAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Session{}, e.ErrSessionNotExist
		}
		log.Fatal(err)
	}
	sess.RefreshToken = refreshToken

	return sess, nil
}

// Errors: -
func (sess *Session) IsExpired() (bool, error) {
	if sess.ExpiresAt.Unix() <= time.Now().Unix() {
		_, err := pgsql.DB.Exec(
			`DELETE FROM sessions WHERE refresh_token=$1`,
			&sess.RefreshToken)
		if err != nil {
			log.Fatal(err)
		}
		return true, nil
	}

	return false, nil
}

// Removes old session
// Errors: -
func ClearSessionsByUserId(user_id int) error {
	stmt, err := pgsql.DB.Prepare(
		`DELETE FROM sessions WHERE user_id=$1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	if _, err := stmt.Exec(user_id); err != nil {
		log.Fatal(err)
	}

	return nil
}
