package userrepo

import (
	"database/sql"
	"fmt"
	"net/mail"
	"time"

	"github.com/google/uuid"

	"github.com/Housiadas/backend-system/internal/core/domain/name"
	"github.com/Housiadas/backend-system/internal/core/domain/role"
	"github.com/Housiadas/backend-system/internal/core/domain/user"
	"github.com/Housiadas/backend-system/pkg/sqldb/dbarray"
)

type userDB struct {
	ID           uuid.UUID      `db:"user_id"`
	Name         string         `db:"name"`
	Email        string         `db:"email"`
	Roles        dbarray.String `db:"roles"`
	PasswordHash []byte         `db:"password_hash"`
	Department   sql.NullString `db:"department"`
	Enabled      bool           `db:"enabled"`
	DateCreated  time.Time      `db:"date_created"`
	DateUpdated  time.Time      `db:"date_updated"`
}

func toUserDB(usr user.User) userDB {
	return userDB{
		ID:           usr.ID,
		Name:         usr.Name.String(),
		Email:        usr.Email.Address,
		Roles:        role.ParseToString(usr.Roles),
		PasswordHash: usr.PasswordHash,
		Department: sql.NullString{
			String: usr.Department.String(),
			Valid:  usr.Department.Valid(),
		},
		Enabled:     usr.Enabled,
		DateCreated: usr.DateCreated.UTC(),
		DateUpdated: usr.DateUpdated.UTC(),
	}
}

func toUserDomain(db userDB) (user.User, error) {
	addr := mail.Address{
		Address: db.Email,
	}

	roles, err := role.ParseMany(db.Roles)
	if err != nil {
		return user.User{}, fmt.Errorf("parse: %w", err)
	}

	nme, err := name.Parse(db.Name)
	if err != nil {
		return user.User{}, fmt.Errorf("parse name: %w", err)
	}

	department, err := name.ParseNull(db.Department.String)
	if err != nil {
		return user.User{}, fmt.Errorf("parse department: %w", err)
	}

	bus := user.User{
		ID:           db.ID,
		Name:         nme,
		Email:        addr,
		Roles:        roles,
		PasswordHash: db.PasswordHash,
		Enabled:      db.Enabled,
		Department:   department,
		DateCreated:  db.DateCreated.In(time.Local),
		DateUpdated:  db.DateUpdated.In(time.Local),
	}

	return bus, nil
}

func toUsersDomain(dbs []userDB) ([]user.User, error) {
	bus := make([]user.User, len(dbs))

	for i, db := range dbs {
		var err error
		bus[i], err = toUserDomain(db)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}
