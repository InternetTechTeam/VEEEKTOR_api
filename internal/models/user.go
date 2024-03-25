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
	Email      string `json:"email"`
	Password   string `json:"password,omitempty"`
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

// Export oriented errors
func (usr *User) Insert() error {
	stmt, err := pgsql.DB.Prepare(
		`INSERT INTO users (
		email, password, name, 
		patronymic, surname, role_id, dep_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}
	if _, err = stmt.Exec(
		&usr.Email, &usr.Password, &usr.Name,
		&usr.Patronymic, &usr.Surname,
		&usr.RoleId, &usr.DepId); err != nil {
		return err
	}
	return nil
}

// Export oriented error message
func (usr *User) PasswordAndLoginValidate() error {
	if len(usr.Password) < PasswordMinLen || len(usr.Password) > PasswordMaxLen {
		msg := fmt.Sprintf(`Password must contain at least %d and no more than %d symbols lenght`, PasswordMinLen, PasswordMaxLen)
		return errors.New(msg)
	}
	return nil
}

// Export oriented errors
func (usr *User) Validate() error {
	if usr.RoleId == 0 || usr.DepId == 0 ||
		usr.Password == "" || usr.Email == "" ||
		usr.Name == "" || usr.Surname == "" {
		return e.ErrMissingFields
	}

	if len(usr.Password) < PasswordMinLen ||
		len(usr.Password) > PasswordMaxLen {
		msg := fmt.Sprintf(`Password must contain at least %d and no more than %d symbols lenght`, PasswordMinLen, PasswordMaxLen)
		return errors.New(msg)
	}

	if len(usr.Name) < FullNameMinLen ||
		len(usr.Name) > FullNameMaxLen {
		msg := fmt.Sprintf(`Name must contain at least %d and no more than %d symbols lenght`, FullNameMinLen, FullNameMaxLen)
		return errors.New(msg)
	}

	if len(usr.Patronymic) < FullNameMinLen ||
		len(usr.Patronymic) > FullNameMaxLen {
		msg := fmt.Sprintf(`Patronymic must contain at least %d and no more than %d symbols lenght`, FullNameMinLen, FullNameMaxLen)
		return errors.New(msg)
	}

	if len(usr.Surname) < FullNameMinLen ||
		len(usr.Surname) > FullNameMaxLen {
		msg := fmt.Sprintf(`Surname must contain at least %d and no more than %d symbols lenght`, FullNameMinLen, FullNameMaxLen)
		return errors.New(msg)
	}

	stmt, err := pgsql.DB.Prepare(
		`SELECT EXISTS(SELECT 1 FROM roles WHERE id=$1)`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	var exists bool
	if stmt.QueryRow(&usr.RoleId).Scan(&exists); !exists {
		return e.ErrRoleNotFound
	}

	_, err = GetDepartmentById(usr.DepId)
	if err != nil {
		return e.ErrDepNotFound
	}

	return nil
}

type SignInInput struct {
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=50"`
}

func (inp *SignInInput) Validate() error {
	if len(inp.Email) < 4 || len(inp.Email) > 64 {
		return errors.New("email not valid")
	}

	if len(inp.Password) < PasswordMinLen ||
		len(inp.Password) > PasswordMaxLen {
		msg := fmt.Sprintf(`Password must contain at least %d and no more than %d symbols lenght`, PasswordMinLen, PasswordMaxLen)
		return errors.New(msg)
	}
	return nil
}

// Export oriented errors
func GetUserByEmailAndPassword(inp SignInInput) (User, error) {
	stmt, err := pgsql.DB.Prepare(
		`SELECT id, email, password, name, 
		patronymic, surname, role_id, dep_id 
		FROM users WHERE email=$1 and password=$2`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	var usr User
	if err := stmt.QueryRow(&inp.Email, &inp.Password).Scan(
		&usr.Id, &usr.Email, &usr.Password,
		&usr.Name, &usr.Patronymic, &usr.Surname,
		&usr.RoleId, &usr.DepId); err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		}
		return usr, e.ErrUserNotFound
	}

	return usr, nil
}

// Export oriented error message
func GetUserById(userId int) (User, error) {
	stmt, err := pgsql.DB.Prepare(`
	SELECT email, password, name, 
	patronymic, surname, role_id, dep_id
	FROM users WHERE id=$1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	var usr User
	usr.Id = userId

	if err := stmt.QueryRow(userId).Scan(
		&usr.Email, &usr.Password, &usr.Name,
		&usr.Patronymic, &usr.Surname,
		&usr.RoleId, &usr.DepId); err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		}
		return usr, e.ErrUserNotFound
	}

	return usr, nil
}
