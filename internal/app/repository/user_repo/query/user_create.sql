INSERT INTO users
(user_id, name, email, password_hash, roles, department, enabled, date_created, date_updated)
VALUES (:user_id, :name, :email, :password_hash, :roles, :department, :enabled, :date_created, :date_updated)