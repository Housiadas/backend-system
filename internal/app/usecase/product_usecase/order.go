package product_usecase

import (
	"github.com/Housiadas/backend-system/internal/core/domain/product"
	"github.com/Housiadas/backend-system/pkg/order"
)

var defaultOrderBy = order.NewBy("product_id", order.ASC)

var orderByFields = map[string]string{
	"product_id": product.OrderByProductID,
	"name":       product.OrderByName,
	"cost":       product.OrderByCost,
	"quantity":   product.OrderByQuantity,
	"user_id":    product.OrderByUserID,
}
