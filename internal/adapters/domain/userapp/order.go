package userapp

import (
	"github.com/Housiadas/backend-system/internal/core/service/userservice"
	"github.com/Housiadas/backend-system/pkg/order"
)

var defaultOrderBy = order.NewBy("user_id", order.ASC)

var orderByFields = map[string]string{
	"user_id": userservice.OrderByID,
	"name":    userservice.OrderByName,
	"email":   userservice.OrderByEmail,
	"roles":   userservice.OrderByRoles,
	"enabled": userservice.OrderByEnabled,
}
