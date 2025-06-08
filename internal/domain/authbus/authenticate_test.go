package authbus_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"github.com/Housiadas/backend-system/internal/domain/authbus"
	"github.com/Housiadas/backend-system/internal/domain/userbus"
	"github.com/Housiadas/backend-system/internal/domain/userbus/stores/userdb"
	"github.com/Housiadas/backend-system/internal/sys/dbtest"
	"github.com/Housiadas/backend-system/internal/sys/types/role"
	"github.com/Housiadas/backend-system/internal/sys/unitest"
)

func Test_Auth(t *testing.T) {

	db := dbtest.NewDatabase(t, "Test_Auth", "file://../../data/migrations")
	sd, err := insertSeedData(db.BusDomain)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	ath := authbus.New(authbus.Config{
		Log:       db.Log,
		DB:        db.DB,
		KeyLookup: &keyStore{},
		Issuer:    "service project",
		Userbus:   userbus.NewBusiness(db.Log, userdb.NewStore(db.Log, db.DB)),
	})

	t.Run("testAdminAuthorization", testAdminAuthorization(ath, sd))
	t.Run("testUserAuthorization", testUserAuthorization(ath, sd))
	t.Run("testUserWithDifferentUUID", testUserWithDifferentUUID(ath, sd))
	t.Run("testUserAdminAuthorization", testUserAdminAuthorization(ath, sd))
	t.Run("testUserRuleAny", testUserRuleAny(ath, sd))
}

func testAdminAuthorization(ath *authbus.Auth, sd unitest.SeedData) func(t *testing.T) {
	f := func(t *testing.T) {
		claims := authbus.Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    ath.Issuer(),
				Subject:   sd.Admins[0].ID.String(),
				ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			},
			Roles: parseRoles([]role.Role{role.Admin}),
		}

		token, err := ath.GenerateToken(claims)
		if err != nil {
			t.Fatalf("Should be able to generate a JWT : %s", err)
		}

		parsedClaims, err := ath.Authenticate(context.Background(), "Bearer "+token)
		if err != nil {
			t.Fatalf("Should be able to authenticate the claims : %s", err)
		}

		userID := uuid.MustParse(claims.Subject)

		err = ath.Authorize(context.Background(), parsedClaims, userID, authbus.RuleAdminOnly)
		if err != nil {
			t.Errorf("Should be able to authorize the Roles.Admin claims : %s", err)
		}

		err = ath.Authorize(context.Background(), parsedClaims, userID, authbus.RuleUserOnly)
		if err == nil {
			t.Error("Should NOT be able to authorize the Roles.User claim")
		}

		err = ath.Authorize(context.Background(), parsedClaims, userID, authbus.RuleAdminOrSubject)
		if err != nil {
			t.Errorf("Should be able to authorize the RuleAdminOrSubject claim with Roles.Admin only : %s", err)
		}
	}

	return f
}

func testUserAuthorization(ath *authbus.Auth, sd unitest.SeedData) func(t *testing.T) {
	f := func(t *testing.T) {
		claims := authbus.Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    ath.Issuer(),
				Subject:   sd.Users[0].ID.String(),
				ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			},
			Roles: parseRoles([]role.Role{role.User}),
		}

		token, err := ath.GenerateToken(claims)
		if err != nil {
			t.Fatalf("Should be able to generate a JWT : %v", err)
		}

		parsedClaims, err := ath.Authenticate(context.Background(), "Bearer "+token)
		if err != nil {
			t.Fatalf("Should be able to authenticate the claims : %s", err)
		}

		userID := uuid.MustParse(claims.Subject)

		err = ath.Authorize(context.Background(), parsedClaims, userID, authbus.RuleUserOnly)
		if err != nil {
			t.Errorf("Should be able to authorize the RuleUserOnly claim with Roles.User only : %s", err)
		}

		err = ath.Authorize(context.Background(), parsedClaims, userID, authbus.RuleAdminOnly)
		if err == nil {
			t.Error("Should NOT be able to authorize the RuleAdminOnly claim with Roles.User only")
		}

		err = ath.Authorize(context.Background(), parsedClaims, userID, authbus.RuleAdminOrSubject)
		if err != nil {
			t.Errorf("Should be able to authorize the RuleAdminOrSubject claim with Roles.User only : %s", err)
		}

		err = ath.Authorize(context.Background(), parsedClaims, userID, authbus.RuleAny)
		if err != nil {
			t.Errorf("Should be able to authorize the RuleAny any claim with Roles.User only : %s", err)
		}
	}

	return f
}

func testUserWithDifferentUUID(ath *authbus.Auth, sd unitest.SeedData) func(t *testing.T) {
	f := func(t *testing.T) {
		claims := authbus.Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    ath.Issuer(),
				Subject:   sd.Users[0].ID.String(),
				ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			},
			Roles: parseRoles([]role.Role{role.User}),
		}

		token, err := ath.GenerateToken(claims)
		if err != nil {
			t.Fatalf("Should be able to generate a JWT : %s", err)
		}

		parsedClaims, err := ath.Authenticate(context.Background(), "Bearer "+token)
		if err != nil {
			t.Fatalf("Should be able to authenticate the claims : %s", err)
		}

		userID := uuid.MustParse("9e979baa-61c9-4b50-81f2-f216d53f5c15")

		err = ath.Authorize(context.Background(), parsedClaims, userID, authbus.RuleAdminOrSubject)
		if err == nil {
			t.Error("Should NOT be able to authorize the RuleAdminOrSubject claim with Roles.User only and different userID")
		}
	}

	return f
}

func testUserAdminAuthorization(ath *authbus.Auth, sd unitest.SeedData) func(t *testing.T) {
	f := func(t *testing.T) {
		claims := authbus.Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    ath.Issuer(),
				Subject:   sd.Admins[0].ID.String(),
				ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			},
			Roles: parseRoles([]role.Role{role.User, role.Admin}),
		}
		userID := uuid.MustParse("9e979baa-61c9-4b50-81f2-f216d53f5c15")

		token, err := ath.GenerateToken(claims)
		if err != nil {
			t.Fatalf("Should be able to generate a JWT : %s", err)
		}

		parsedClaims, err := ath.Authenticate(context.Background(), "Bearer "+token)
		if err != nil {
			t.Fatalf("Should be able to authenticate the claims : %s", err)
		}

		err = ath.Authorize(context.Background(), parsedClaims, userID, authbus.RuleAny)
		if err != nil {
			t.Errorf("Should be able to authorize the RuleAny any claim with Roles.User and Roles.Admin : %s", err)
		}
	}

	return f
}

func testUserRuleAny(ath *authbus.Auth, sd unitest.SeedData) func(t *testing.T) {
	f := func(t *testing.T) {
		claims := authbus.Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    ath.Issuer(),
				Subject:   sd.Users[0].ID.String(),
				ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			},
			Roles: parseRoles([]role.Role{role.User}),
		}
		userID := uuid.MustParse("9e979baa-61c9-4b50-81f2-f216d53f5c15")

		token, err := ath.GenerateToken(claims)
		if err != nil {
			t.Fatalf("Should be able to generate a JWT : %s", err)
		}

		parsedClaims, err := ath.Authenticate(context.Background(), "Bearer "+token)
		if err != nil {
			t.Fatalf("Should be able to authenticate the claims : %s", err)
		}

		err = ath.Authorize(context.Background(), parsedClaims, userID, authbus.RuleAny)
		if err != nil {
			t.Errorf("Should be able to authorize the RuleAny any claim with Roles.User only : %s", err)
		}
	}

	return f
}

// =============================================================================

func insertSeedData(busDomain dbtest.BusDomain) (unitest.SeedData, error) {
	ctx := context.Background()

	usrs, err := userbus.TestSeedUsers(ctx, 2, role.Admin, busDomain.User)
	if err != nil {
		return unitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	tu1 := unitest.User{
		User: usrs[0],
	}

	tu2 := unitest.User{
		User: usrs[1],
	}

	// -------------------------------------------------------------------------

	usrs, err = userbus.TestSeedUsers(ctx, 2, role.User, busDomain.User)
	if err != nil {
		return unitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	tu3 := unitest.User{
		User: usrs[0],
	}

	tu4 := unitest.User{
		User: usrs[1],
	}

	// -------------------------------------------------------------------------

	sd := unitest.SeedData{
		Users:  []unitest.User{tu3, tu4},
		Admins: []unitest.User{tu1, tu2},
	}

	return sd, nil
}

// =============================================================================

type keyStore struct{}

func (ks *keyStore) PrivateKey() (string, error) {
	return privateKeyPEM, nil
}

func (ks *keyStore) PublicKey() (string, error) {
	return publicKeyPEM, nil
}

func parseRoles(roles []role.Role) []string {
	appRoles := make([]string, len(roles))
	for i, r := range roles {
		appRoles[i] = r.String()
	}
	return appRoles
}

const (
	privateKeyPEM = `-----BEGIN PRIVATE KEY-----
MIIEpQIBAAKCAQEAvMAHb0IoLvoYuW2kA+LTmnk+hfnBq1eYIh4CT/rMPCxgtzjq
U0guQOMnLg69ydyA5uu37v6rbS1+stuBTEiMQl/bxAhgLkGrUhgpZ10Bt6GzSEgw
QNloZoGaxe4p20wMPpT4kcMKNHkQds3uONNcLxPUmfjbbH64g+seg28pbgQPwKFK
tF7bIsOBgz0g5Ptn5mrkdzqMPUSy9k9VCu+R42LH9c75JsRzz4FeN+VzwMAL6yQn
ZvOi7/zOgNyxeVia8XVKykrnhgcpiOn5oaLRBzQGN00Z7TuBRIfDJWU21qQN4Cq7
keZmMP4gqCVWjYneK4bzrG/+H2w9BJ2TsmMGvwIDAQABAoIBAFQmQKpHkmavNYql
6POaksBRwaA1YzSijr7XJizGIXvKRSwqgb2zdnuTSgpspAx09Dr/aDdy7rZ0DAJt
fk2mInINDottOIQm3txwzTS58GQQAT/+fxTKWJMqwPfxYFPWqbbU76T8kXYna0Gs
OcK36GdMrgIfQqQyMs0Na8MpMg1LmkAxuqnFCXS/NMyKl9jInaaTS+Kz+BSzUMGQ
zebfLFsf2N7sLZuimt9zlRG30JJTfBlB04xsYMo734usA2ITe8U0XqG6Og0qc6ev
6lsoM8hpvEUsQLcjQQ5up7xx3S2stZJ8o0X8GEX5qUMaomil8mZ7X5xOlEqf7p+v
lXQ46cECgYEA2lbZQON6l3ZV9PCn9j1rEGaXio3SrAdTyWK3D1HF+/lEjClhMkfC
XrECOZYj+fiI9n+YpSog+tTDF7FTLf7VP21d2gnhQN6KAXUnLIypzXxodcC6h+8M
ZGJh/EydLvC7nPNoaXx96bohxzS8hrOlOlkCbr+8gPYKf8qkbe7HyxECgYEA3U6e
x9g4FfTvI5MGrhp2BIzoRSn7HlNQzjJ71iMHmM2kBm7TsER8Co1PmPDrP8K/UyGU
Q25usTsPSrHtKQEV6EsWKaP/6p2Q82sDkT9bZlV+OjRvOfpdO5rP6Q95vUmMGWJ/
S6oimbXXL8p3gDafw3vC1PCAhoaxMnGyKuZwlM8CgYEAixT1sXr2dZMg8DV4mMfI
8pqXf+AVyhWkzsz+FVkeyAKiIrKdQp0peI5C/5HfevVRscvX3aY3efCcEfSYKt2A
07WEKkdO4LahrIoHGT7FT6snE5NgfwTMnQl6p2/aVLNun20CHuf5gTBbIf069odr
Af7/KLMkjfWs/HiGQ6zuQjECgYEAv+DIvlDz3+Wr6dYyNoXuyWc6g60wc0ydhQo0
YKeikJPLoWA53lyih6uZ1escrP23UOaOXCDFjJi+W28FR0YProZbwuLUoqDW6pZg
U3DxWDrL5L9NqKEwcNt7ZIDsdnfsJp5F7F6o/UiyOFd9YQb7YkxN0r5rUTg7Lpdx
eMyv0/UCgYEAhX9MPzmTO4+N8naGFof1o8YP97pZj0HkEvM0hTaeAQFKJiwX5ijQ
xumKGh//G0AYsjqP02ItzOm2mWnbI3FrNlKmGFvR6VxIZMOyXvpLofHucjJ5SWli
eYjPklKcXaMftt1FVO4n+EKj1k1+Tv14nytq/J5WN+r4FBlNEYj/6vg=
-----END PRIVATE KEY-----
`
	publicKeyPEM = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvMAHb0IoLvoYuW2kA+LT
mnk+hfnBq1eYIh4CT/rMPCxgtzjqU0guQOMnLg69ydyA5uu37v6rbS1+stuBTEiM
Ql/bxAhgLkGrUhgpZ10Bt6GzSEgwQNloZoGaxe4p20wMPpT4kcMKNHkQds3uONNc
LxPUmfjbbH64g+seg28pbgQPwKFKtF7bIsOBgz0g5Ptn5mrkdzqMPUSy9k9VCu+R
42LH9c75JsRzz4FeN+VzwMAL6yQnZvOi7/zOgNyxeVia8XVKykrnhgcpiOn5oaLR
BzQGN00Z7TuBRIfDJWU21qQN4Cq7keZmMP4gqCVWjYneK4bzrG/+H2w9BJ2TsmMG
vwIDAQAB
-----END PUBLIC KEY-----
`
)
