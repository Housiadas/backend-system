package productapp

import (
	"github.com/Housiadas/backend-system/internal/core/service/productbus"
	"github.com/Housiadas/backend-system/pkg/order"
)

var defaultOrderBy = order.NewBy("product_id", order.ASC)

var orderByFields = map[string]string{
	"product_id": productbus.OrderByProductID,
	"name":       productbus.OrderByName,
	"cost":       productbus.OrderByCost,
	"quantity":   productbus.OrderByQuantity,
	"user_id":    productbus.OrderByUserID,
}
