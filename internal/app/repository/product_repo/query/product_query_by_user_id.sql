SELECT product_id,
       user_id,
       name,
       cost,
       quantity,
       date_created,
       date_updated
FROM products
WHERE user_id = :user_id