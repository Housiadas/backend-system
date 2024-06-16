-- Description: Add products view.
CREATE OR REPLACE VIEW view_user_products AS
SELECT p.product_id,
       p.user_id,
       p.name,
       p.cost,
       p.quantity,
       p.date_created,
       p.date_updated,
       u.name AS user_name
FROM products AS p
         JOIN
     users AS u ON u.user_id = p.user_id
