UPDATE
    products
SET "name"         = :name,
    "cost"         = :cost,
    "quantity"     = :quantity,
    "date_updated" = :date_updated
WHERE product_id = :product_id