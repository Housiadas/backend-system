CREATE TABLE "users"
(
    "username"            varchar PRIMARY KEY,
    "role"                varchar        NOT NULL DEFAULT 'depositor',
    "hashed_password"     varchar        NOT NULL,
    "full_name"           varchar        NOT NULL,
    "email"               varchar UNIQUE NOT NULL,
    "is_email_verified"   bool           NOT NULL DEFAULT false,
    "password_changed_at" timestamptz    NOT NULL DEFAULT '0001-01-01',
    "created_at"          timestamptz    NOT NULL DEFAULT (now())
);

CREATE TABLE "verify_emails"
(
    "id"          bigserial PRIMARY KEY,
    "username"    varchar     NOT NULL REFERENCES "users" ("username") ON DELETE CASCADE,
    "email"       varchar     NOT NULL,
    "secret_code" varchar     NOT NULL,
    "is_used"     bool        NOT NULL DEFAULT false,
    "created_at"  timestamptz NOT NULL DEFAULT (now()),
    "expired_at"  timestamptz NOT NULL DEFAULT (now() + interval '15 minutes')
);

CREATE TABLE "accounts"
(
    "id"         bigserial PRIMARY KEY,
    "owner"      varchar     NOT NULL REFERENCES "users" ("username") ON DELETE CASCADE,
    "balance"    bigint      NOT NULL,
    "currency"   varchar     NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "entries"
(
    "id"         bigserial PRIMARY KEY,
    "account_id" bigint      NOT NULL REFERENCES "accounts" ("id") ON DELETE CASCADE,
    "amount"     bigint      NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

COMMENT ON COLUMN "entries"."amount" IS 'can be negative or positive';

CREATE TABLE "transfers"
(
    "id"              bigserial PRIMARY KEY,
    "from_account_id" bigint      NOT NULL REFERENCES "accounts" ("id") ON DELETE CASCADE,
    "to_account_id"   bigint      NOT NULL REFERENCES "accounts" ("id") ON DELETE CASCADE,
    "amount"          bigint      NOT NULL,
    "created_at"      timestamptz NOT NULL DEFAULT (now())
);

COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';

CREATE TABLE "sessions"
(
    "id"            uuid PRIMARY KEY,
    "username"      varchar     NOT NULL REFERENCES "users" ("username") ON DELETE CASCADE,
    "refresh_token" varchar     NOT NULL,
    "user_agent"    varchar     NOT NULL,
    "client_ip"     varchar     NOT NULL,
    "is_blocked"    boolean     NOT NULL DEFAULT false,
    "expires_at"    timestamptz NOT NULL,
    "created_at"    timestamptz NOT NULL DEFAULT (now())
);