package repository

import (
	"database/sql"
	"errors"

	"github.com/go-openapi/strfmt"
	"github.com/sylph4/entain-task/storage/postgres"
	"gopkg.in/gorp.v1"
)

type IUserRepository interface {
	Update(user *postgres.User) error
	SelectUserByID(id strfmt.UUID4) (*postgres.User, error)
}

type UserRepository struct {
	db *gorp.DbMap
}

func NewUserRepository(conn *gorp.DbMap) *UserRepository {
	return &UserRepository{db: conn}
}

func (s *UserRepository) SelectUserByID(tx *gorp.Transaction, id strfmt.UUID4) (*postgres.User, error) {
	user := &postgres.User{}
	if err := tx.SelectOne(&user, `
		SELECT
			*
		FROM
			users
		WHERE
			id = $1
	`,
		id,
	); err != nil {
		if err == sql.ErrNoRows {
			// user not found
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}

func (s *UserRepository) Update(tx *gorp.Transaction, user postgres.User) error {
	_, err := tx.Update(&user)
	if err != nil {
		return errors.New("updating user failed")
	}

	return nil
}
