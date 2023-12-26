package db

import (
	"context"
	"github.com/Housiadas/simple-banking-system/foundation/random"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := random.HashPassword(random.RandomString(6))
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
