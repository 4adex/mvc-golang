package models

import (
	"fmt"
	"github.com/4adex/mvc-golang/pkg/types"
)

func GetUser(field string) (types.User, error) {
	var User types.User
	db, err := Connection()
	if err != nil {
		return User, err
	}
	query := "SELECT * FROM users WHERE username = ? or email = ?"
	rows, err := db.Query(query, field, field)
	if err != nil {
		return User, err
	}
	if rows.Next() {
		err = rows.Scan(&User.ID, &User.Username, &User.Password, &User.Email, &User.Role, &User.Salt, &User.RequestStatus)
		if err != nil {
			return User, err
		}
	} else {
		return User, fmt.Errorf("no user found")
	}
	return User, nil
}

func CreateUser(user types.User) error {
	db, err := Connection()
	if err != nil {
		return err
	}
	query := "INSERT INTO users (username, password, email, role, salt, request_status) VALUES (?, ?, ?, ?, ?, ?)"
	_, err = db.Exec(query, user.Username, user.Password, user.Email, user.Role, user.Salt, user.RequestStatus)
	if err != nil {
		return err
	}
	return nil
}

func IsUsersTableEmpty() (bool, error) {
	db, err := Connection()
	if err != nil {
		return false, err
	}

	var count int
	query := "SELECT COUNT(*) FROM users"
	err = db.QueryRow(query).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}