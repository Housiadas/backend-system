UPDATE
    users
SET "name"          = :name,
    "email"         = :email,
    "roles"         = :roles,
    "password_hash" = :password_hash,
    "department"    = :department,
    "enabled"       = :enabled,
    "date_updated"  = :date_updated
WHERE user_id = :user_id