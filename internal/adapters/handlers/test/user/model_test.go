package user_test

import (
	"time"

	"github.com/Housiadas/backend-system/internal/adapters/domain/userapp"
	"github.com/Housiadas/backend-system/internal/core/domain/role"
	"github.com/Housiadas/backend-system/internal/core/service/userservice"
)

func toAppUser(bus userservice.User) userapp.User {
	return userapp.User{
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

func toAppUsers(users []userservice.User) []userapp.User {
	items := make([]userapp.User, len(users))
	for i, usr := range users {
		items[i] = toAppUser(usr)
	}

	return items
}

func toAppUserPtr(bus userservice.User) *userapp.User {
	appUsr := toAppUser(bus)
	return &appUsr
}
