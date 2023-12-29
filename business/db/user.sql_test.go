package db

import (
	"context"
	"testing"

	"github.com/Housiadas/simple-banking-system/foundation/password"
	"github.com/Housiadas/simple-banking-system/foundation/random"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := password.HashPassword(random.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       random.RandomUsername(),
		HashedPassword: hashedPassword,
		FullName:       random.RandomUsername(),
		Email:          random.RandomEmail(),
	}

	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}
