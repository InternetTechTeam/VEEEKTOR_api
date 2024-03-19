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

// User id and role_id must be valid
func StoreSession(user_id int, role_id int) (TokenResponse, error) {
	stmt, err := pgsql.DB.Prepare(
		`INSERT INTO sessions (user_id, refresh_token, expires_at)
		VALUES ($1, $2, $3)`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	var resp TokenResponse
	resp.AccessToken, _ = GenerateAccessToken(user_id, role_id)
	resp.RefreshToken, _ = GenerateRefreshToken()

	CheckSessionsCount(user_id)

	if _, err := stmt.Exec(
		&user_id, &resp.RefreshToken,
		time.Now().Add(RefreshTokenLifeTime)); err != nil {
		log.Print(err)
		return TokenResponse{}, err
	}

	return resp, nil
}

func DeleteSessionByRT(refreshToken string) {
	stmt, err := pgsql.DB.Prepare(
		`DELETE FROM sessions WHERE refresh_token=$1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	if _, err := stmt.Exec(&refreshToken); err != nil {
		log.Fatal(err)
	}
}

func CheckSessionsCount(user_id int) {
	getStmt, err := pgsql.DB.Prepare(
		`SELECT COUNT(*) FROM sessions WHERE user_id=$1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}
	var count int
	if getStmt.QueryRow(&user_id).Scan(&count); count <= 5 {
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

func GetSessionByRefreshToken(refreshToken string) (Session, error) {
	stmt, err := pgsql.DB.Prepare(
		`SELECT id, user_id, expires_at FROM 
		sessions WHERE refresh_token=$1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	var sess Session
	if err := stmt.QueryRow(&refreshToken).Scan(
		&sess.Id, &sess.UserId, &sess.ExpiresAt); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return Session{}, err
		}
		return Session{}, e.ErrSessionNotExist
	}
	return sess, nil
}

func (sess *Session) IsExpired() (bool, error) {
	if sess.ExpiresAt.Unix() <= time.Now().Unix() {
		stmt, err := pgsql.DB.Prepare(
			`DELETE FROM sessions WHERE refresh_token=$1`)
		if err != nil {
			log.Fatal(e.ErrCantPrepareDbStmt)
		}

		if _, err = stmt.Exec(sess.RefreshToken); err != nil {
			return true, err
		}
	}

	return false, nil
}

// Removes old session
func ClearSessionsByUserId(user_id int) error {
	stmt, err := pgsql.DB.Prepare(
		`DELETE FROM sessions WHERE user_id=$1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	if _, err := stmt.Exec(user_id); err != nil {
		log.Print(err)
		return err
	}

	return nil
}
