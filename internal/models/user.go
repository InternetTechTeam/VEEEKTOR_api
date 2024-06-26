package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"VEEEKTOR_api/pkg/database/pgsql"
	e "VEEEKTOR_api/pkg/errors"
)

type User struct {
	Id         int    `json:"id"`
	Email      string `json:"email,omitempty"`
	Password   string `json:"password,omitempty"`
	GroupId    int    `json:"group_id"`
	Name       string `json:"name"`
	Patronymic string `json:"patronymic"`
	Surname    string `json:"surname"`
	RoleId     int    `json:"role_id"`
	DepId      int    `json:"dep_id"`
}

const (
	PasswordMinLen = 8
	PasswordMaxLen = 50
	FullNameMinLen = 2
	FullNameMaxLen = 30
)

// Errors: ErrUserNotFound
func GetUserById(userId int) (User, error) {
	stmt, err := pgsql.DB.Prepare(`
	SELECT email, password, group_id, name, 
	patronymic, surname, role_id, dep_id
	FROM users WHERE id=$1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	var usr User
	usr.Id = userId

	if err := stmt.QueryRow(userId).Scan(
		&usr.Email, &usr.Password, &usr.GroupId,
		&usr.Name, &usr.Patronymic, &usr.Surname,
		&usr.RoleId, &usr.DepId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return usr, e.ErrUserNotFound
		}
		log.Fatal(err)
	}

	return usr, nil
}

// Errors: ErrUserNotFound
func GetUserByEmailAndPassword(inp SignInInput) (User, error) {
	stmt, err := pgsql.DB.Prepare(
		`SELECT id, email, password, group_id, name, 
		patronymic, surname, role_id, dep_id 
		FROM users WHERE email=$1 and password=$2`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	var usr User
	if err := stmt.QueryRow(&inp.Email, &inp.Password).Scan(
		&usr.Id, &usr.Email, &usr.Password,
		&usr.GroupId, &usr.Name, &usr.Patronymic,
		&usr.Surname, &usr.RoleId, &usr.DepId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return usr, e.ErrUserNotFound
		}
		log.Fatal(err)
	}

	return usr, nil
}

// Errors: message, ErrMissingFields, ErrRoleNotFound,
// ErrDepNotFound, ErrUserExists
func (usr *User) Validate() error {
	if usr.RoleId == 0 || usr.DepId == 0 ||
		usr.Password == "" || usr.Email == "" ||
		usr.GroupId == 0 || usr.Name == "" ||
		usr.Surname == "" {
		return e.ErrMissingFields
	}

	if len(usr.Password) < PasswordMinLen ||
		len(usr.Password) > PasswordMaxLen {
		msg := fmt.Sprintf(`password must contain at least %d and no more than %d symbols lenght`, PasswordMinLen, PasswordMaxLen)
		return errors.New(msg)
	}

	if len(usr.Name) < FullNameMinLen ||
		len(usr.Name) > FullNameMaxLen {
		msg := fmt.Sprintf(`name must contain at least %d and no more than %d symbols lenght`, FullNameMinLen, FullNameMaxLen)
		return errors.New(msg)
	}

	if len(usr.Patronymic) < FullNameMinLen ||
		len(usr.Patronymic) > FullNameMaxLen {
		msg := fmt.Sprintf(`patronymic must contain at least %d and no more than %d symbols lenght`, FullNameMinLen, FullNameMaxLen)
		return errors.New(msg)
	}

	if len(usr.Surname) < FullNameMinLen ||
		len(usr.Surname) > FullNameMaxLen {
		msg := fmt.Sprintf(`surname must contain at least %d and no more than %d symbols lenght`, FullNameMinLen, FullNameMaxLen)
		return errors.New(msg)
	}

	var exists bool
	err := pgsql.DB.QueryRow(
		`SELECT 1 FROM roles WHERE id=$1`,
		&usr.RoleId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e.ErrRoleNotFound
		}
		log.Fatal(err)
	}

	err = pgsql.DB.QueryRow(
		`SELECT 1 FROM groups WHERE id=$1`,
		&usr.GroupId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e.ErrGroupNotExist
		}
		log.Fatal(err)
	}

	err = pgsql.DB.QueryRow(
		`SELECT 1 FROM departments WHERE id=$1`,
		&usr.DepId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e.ErrDepNotFound
		}
		log.Fatal(err)
	}

	err = pgsql.DB.QueryRow(
		`SELECT 1 FROM users WHERE email=$1`,
		&usr.Email).Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Fatal(err)
	} else if err == nil {
		return e.ErrUserExist
	}

	return nil
}

// Errors: message, ErrMissingFields, ErrRoleNotFound,
// ErrDepNotFound
func (usr *User) Insert() error {
	if err := usr.Validate(); err != nil {
		return err
	}

	stmt, err := pgsql.DB.Prepare(
		`INSERT INTO users (
		email, password, group_id, name, 
		patronymic, surname, role_id, dep_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}
	if _, err = stmt.Exec(
		&usr.Email, &usr.Password, &usr.GroupId,
		&usr.Name, &usr.Patronymic, &usr.Surname,
		&usr.RoleId, &usr.DepId); err != nil {
		log.Fatal(err)
	}
	return nil
}

type SignInInput struct {
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=50"`
}

// Errors: message
func (inp *SignInInput) Validate() error {
	if len(inp.Email) < 4 || len(inp.Email) > 64 {
		return errors.New("email not valid")
	}

	if len(inp.Password) < PasswordMinLen ||
		len(inp.Password) > PasswordMaxLen {
		msg := fmt.Sprintf(`password must contain at least %d and no more than %d symbols lenght`, PasswordMinLen, PasswordMaxLen)
		return errors.New(msg)
	}
	return nil
}
