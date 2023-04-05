package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `insert into users(name, email, hashed_password, created) 
	values ($1, $2, $3, now() at time zone 'utc')`
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		var postgresqlErr *pq.Error
		if errors.As(err, &postgresqlErr) {
			if postgresqlErr.Code == "23505" &&
				strings.Contains(postgresqlErr.Message, "users_us_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
