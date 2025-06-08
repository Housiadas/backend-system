package commands

import (
	"context"
	"fmt"
	"net/mail"
	"time"

	"github.com/Housiadas/backend-system/internal/data/sqldb"
	"github.com/Housiadas/backend-system/internal/domain/userbus"
	"github.com/Housiadas/backend-system/internal/domain/userbus/stores/userdb"
	namePck "github.com/Housiadas/backend-system/internal/sys/types/name"
	"github.com/Housiadas/backend-system/internal/sys/types/role"
)

// UserAdd adds new users into the database.
func (cmd *Command) UserAdd(name, email, password string) error {
	if name == "" || email == "" || password == "" {
		fmt.Println("help: useradd <name> <email> <password>")
		return ErrHelp
	}

	db, err := sqldb.Open(cmd.DB)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userBus := userbus.NewBusiness(cmd.Log, userdb.NewStore(cmd.Log, db))

	addr, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("parsing email: %w", err)
	}

	nu := userbus.NewUser{
		Name:     namePck.MustParse(name),
		Email:    *addr,
		Password: password,
		Roles:    []role.Role{role.Admin, role.User},
	}

	usr, err := userBus.Create(ctx, nu)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	fmt.Println("user id:", usr.ID)
	return nil
}
