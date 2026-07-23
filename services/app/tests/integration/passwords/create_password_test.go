package passwords_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	passwords_api "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/passwords/api"
	passwords_postgres_repository "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/passwords/repository/postgres"
	passwords_service "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/passwords/service"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	USERID  = 1
	TIMES   = 3
	MEMORY  = 64 * 1024
	THREADS = 2
	KEYLEN  = 32
)

func TestCreatePassword(t *testing.T) {
	ctx := context.Background()

	passwordHasher := passwords_service.NewArgon2IDHasher(
		TIMES,
		MEMORY,
		THREADS,
		KEYLEN,
	)

	db, connString := SetupTestDB(t, ctx)
	RunMigrations(t, db, connString)

	_, err := db.Exec(
		ctx,
		`
		INSERT INTO app.users (id, version, name, surname, username, birth_date, description, email, phone_number)
		VALUES (1, 1, 'Alex', 'Bobkov', 'alexey', '2005-03-19', 'I am a developer', 'alexeybobkovdev@outlook.com', '+79200675223');
	`,
	)
	require.NoError(t, err)

	repository := passwords_postgres_repository.NewPasswordsPostgresRepository(db)
	service := passwords_service.NewPasswordService(
		repository,
		passwordHasher,
	)
	passwordAPI := passwords_api.NewPasswordAPI(service)

	passwordToHash := []byte("Password-TO-HASH-2026")

	password, err := passwordAPI.CreatePassword(
		ctx,
		USERID,
		passwordToHash,
	)
	require.NoError(t, err, "failed to create password")
	require.NotEqual(t, passwordToHash, password.Hash, "passwordToHash must not be equal to hashed password")
}

func SetupTestDB(t *testing.T, ctx context.Context) (*pgxpool.Pool, string) {
	t.Helper()

	containerDB, err := postgres.Run(
		ctx,
		"postgres:18.4-alpine",
		postgres.WithDatabase("tripmate_test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
		),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = containerDB.Terminate(ctx) })

	connString, err := containerDB.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connString)
	require.NoError(t, err)

	t.Cleanup(func() { pool.Close() })

	return pool, connString
}

func RunMigrations(t *testing.T, db *pgxpool.Pool, connString string) {
	t.Helper()

	migrationsRoot := getMigrationsRoot()

	m, err := migrate.New(
		"file://"+migrationsRoot,
		connString,
	)
	if err != nil {
		require.NoError(t, err, "create migrate instance")
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		require.NoError(t, err, "run migrations")
	}
}

func getProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)

		if parent == dir {
			panic("could not find root directory")
		}

		dir = parent
	}
}

func getMigrationsRoot() string {
	projectRoot := getProjectRoot()
	return filepath.Join(projectRoot, "migrations")
}
