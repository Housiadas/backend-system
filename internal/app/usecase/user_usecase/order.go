package user_usecase

import (
	"github.com/Housiadas/backend-system/internal/core/domain/user"
	"github.com/Housiadas/backend-system/pkg/order"
)

var defaultOrderBy = order.NewBy("user_id", order.ASC)

var orderByFields = map[string]string{
	"user_id": user.OrderByID,
	"name":    user.OrderByName,
	"email":   user.OrderByEmail,
	"roles":   user.OrderByRoles,
	"enabled": user.OrderByEnabled,
}
