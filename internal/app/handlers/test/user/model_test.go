package user_test

import (
	"time"

	"github.com/Housiadas/backend-system/internal/app/usecase/user_usecase"
	"github.com/Housiadas/backend-system/internal/core/domain/role"
	"github.com/Housiadas/backend-system/internal/core/domain/user"
)

func toAppUser(bus user.User) user_usecase.User {
	return user_usecase.User{
		ID:           bus.ID.String(),
		Name:         bus.Name.String(),
		Email:        bus.Email.Address,
		Roles:        role.ParseToString(bus.Roles),
		PasswordHash: nil, // This field is not marshalled.
		Department:   bus.Department.String(),
		Enabled:      bus.Enabled,
		DateCreated:  bus.DateCreated.Format(time.RFC3339),
		DateUpdated:  bus.DateUpdated.Format(time.RFC3339),
	}
}

func toAppUsers(users []user.User) []user_usecase.User {
	items := make([]user_usecase.User, len(users))
	for i, usr := range users {
		items[i] = toAppUser(usr)
	}

	return items
}

func toAppUserPtr(bus user.User) *user_usecase.User {
	appUsr := toAppUser(bus)
	return &appUsr
}
