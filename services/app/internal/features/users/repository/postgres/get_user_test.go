package users_postgres_repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	pool_mocks "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/mocks"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

func TestUsersRepository_GetUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pool_mocks.NewMockPool(ctrl)
	repo := NewUsersRepository(mockPool)

	ctx := context.Background()
	userID := 1

	// Create a test user
	testUser := testutils.UserTest()

	// Mock the OpTimeout method
	mockPool.
		EXPECT().
		OpTimeout().
		Return(15 * time.Second)

	// Mock the QueryRow method to return the test user
	mockPool.
		EXPECT().
		QueryRow(gomock.Any(), getUserQuery, userID).
		Return(&mockRow{
			err: nil,
		})

	retrievedUser, err := repo.GetUser(ctx, userID)
	require.NoError(t, err)
	require.NotNil(t, retrievedUser)

	assert.Equal(t, testUser.ID, retrievedUser.ID)
	assert.Equal(t, testUser.Name, retrievedUser.Name)
	assert.Equal(t, testUser.Email, retrievedUser.Email)
}

func TestUsersRepository_GetUser_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pool_mocks.NewMockPool(ctrl)
	repo := NewUsersRepository(mockPool)

	ctx := context.Background()
	userID := 999

	// Mock the OpTimeout method
	mockPool.
		EXPECT().
		OpTimeout().
		Return(15 * time.Second)

	// Mock the QueryRow method to return no rows
	mockPool.
		EXPECT().
		QueryRow(gomock.Any(), getUserQuery, userID).
		Return(&mockRow{err: pgx.ErrNoRows})

	retrievedUser, err := repo.GetUser(ctx, userID)
	require.Error(t, err)
	assert.Nil(t, retrievedUser)

	// Check that the error is wrapped correctly
	assert.ErrorContains(t, err, "get user")
	assert.ErrorContains(t, err, "no rows in result set")
}

func TestUsersRepository_GetUser_Deleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pool_mocks.NewMockPool(ctrl)
	repo := NewUsersRepository(mockPool)

	ctx := context.Background()
	userID := 1

	// Mock the OpTimeout method
	mockPool.
		EXPECT().
		OpTimeout().
		Return(15 * time.Second)

	// Mock the QueryRow method to return no rows (because deleted_at is not NULL)
	mockPool.
		EXPECT().
		QueryRow(gomock.Any(), getUserQuery, userID).
		Return(&mockRow{err: pgx.ErrNoRows})

	retrievedUser, err := repo.GetUser(ctx, userID)
	require.Error(t, err)
	assert.Nil(t, retrievedUser)

	// Check that the error is wrapped correctly
	assert.ErrorContains(t, err, "get user")
	assert.ErrorContains(t, err, "no rows in result set")
}

func TestUsersRepository_GetUser_DatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pool_mocks.NewMockPool(ctrl)
	repo := NewUsersRepository(mockPool)

	ctx := context.Background()
	userID := 1

	// Mock the OpTimeout method
	mockPool.
		EXPECT().
		OpTimeout().
		Return(15 * time.Second)

	// Mock the QueryRow method to return a database error
	mockPool.
		EXPECT().
		QueryRow(gomock.Any(), getUserQuery, userID).
		Return(&mockRow{err: errors.New("database error")})

	retrievedUser, err := repo.GetUser(ctx, userID)
	require.Error(t, err)
	assert.Nil(t, retrievedUser)

	// Check that the error is wrapped correctly
	assert.ErrorContains(t, err, "get user")
}

func TestUsersRepository_GetUser_SpecificDatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pool_mocks.NewMockPool(ctrl)
	repo := NewUsersRepository(mockPool)

	ctx := context.Background()
	userID := 1

	// Mock the OpTimeout method
	mockPool.
		EXPECT().
		OpTimeout().
		Return(15 * time.Second)

	// Mock the QueryRow method to return a specific database error
	dbError := errors.New("connection timeout")
	mockPool.
		EXPECT().
		QueryRow(gomock.Any(), getUserQuery, userID).
		Return(&mockRow{err: dbError})

	retrievedUser, err := repo.GetUser(ctx, userID)
	require.Error(t, err)
	assert.Nil(t, retrievedUser)

	// Check that the error is wrapped correctly
	assert.ErrorContains(t, err, "get user")
	assert.ErrorContains(t, err, "connection timeout")
}

func TestUsersRepository_GetUser_ContextTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pool_mocks.NewMockPool(ctrl)
	repo := NewUsersRepository(mockPool)

	ctx := context.Background()
	userID := 1

	// Mock the OpTimeout method to return a very short timeout
	mockPool.
		EXPECT().
		OpTimeout().
		Return(1 * time.Nanosecond)

	// Mock the QueryRow method to simulate a timeout
	mockPool.
		EXPECT().
		QueryRow(gomock.Any(), getUserQuery, userID).
		Return(&mockRow{err: context.DeadlineExceeded})

	retrievedUser, err := repo.GetUser(ctx, userID)
	require.Error(t, err)
	assert.Nil(t, retrievedUser)

	// Check that the error is wrapped correctly
	assert.ErrorContains(t, err, "get user")
	assert.ErrorContains(t, err, "deadline exceeded")
}

func TestUsersRepository_Get_TransactionCommitError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pool_mocks.NewMockPool(ctrl)
	repo := NewUsersRepository(mockPool)

	ctx := context.Background()
	userID := 1

	// Mock the database to return a specific error
	mockPool.
		EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.
		EXPECT().
		QueryRow(gomock.Any(), getUserQuery, userID).
		Return(&mockRow{err: errors.New("database error")})

	retrievedUser, err := repo.GetUser(ctx, userID)
	require.Error(t, err)
	assert.Nil(t, retrievedUser)

	// Check that the error is wrapped correctly
	assert.ErrorContains(t, err, "get user")
	assert.ErrorContains(t, err, "database error")
}

// Mock row for testing
type mockRow struct {
	err error
}

func (m *mockRow) Scan(dest ...interface{}) error {
	if m.err != nil {
		return m.err
	}

	// Create a test user with the expected values
	testUser := testutils.UserTest()

	// Simulate scanning the user into the destination pointers
	if len(dest) > 0 {
		*dest[0].(*int) = testUser.ID
	}
	if len(dest) > 1 {
		*dest[1].(*int) = testUser.Version
	}
	if len(dest) > 2 {
		*dest[2].(*string) = testUser.Name
	}
	if len(dest) > 3 {
		*dest[3].(*string) = testUser.Surname
	}
	if len(dest) > 4 {
		*dest[4].(*string) = testUser.Username
	}
	if len(dest) > 5 {
		*dest[5].(*time.Time) = testUser.BirthDate
	}
	if len(dest) > 6 {
		*dest[6].(*string) = testUser.Description
	}
	if len(dest) > 7 {
		*dest[7].(*string) = testUser.Email
	}
	if len(dest) > 8 {
		*dest[8].(*string) = testUser.PhoneNumber
	}
	if len(dest) > 9 {
		*dest[9].(*time.Time) = testUser.CreatedAt
	}
	if len(dest) > 10 {
		*dest[10].(*time.Time) = testUser.UpdatedAt
	}
	if len(dest) > 11 {
		if testUser.DeletedAt != nil {
			*dest[11].(*time.Time) = *testUser.DeletedAt
		} else {
			dest[11] = nil
		}
	}

	return nil
}

func TestUsersRepository_GetUser_ErrorHandling(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pool_mocks.NewMockPool(ctrl)
	repo := NewUsersRepository(mockPool)

	ctx := context.Background()
	userID := 1

	// Mock the database to return a specific error
	mockPool.
		EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.
		EXPECT().
		QueryRow(gomock.Any(), getUserQuery, userID).
		Return(&mockRow{err: errors.New("database error")})

	retrievedUser, err := repo.GetUser(ctx, userID)
	require.Error(t, err)
	assert.Nil(t, retrievedUser)

	// Check that the error is wrapped correctly
	assert.ErrorContains(t, err, "get user")
	assert.ErrorContains(t, err, "database error")
}
