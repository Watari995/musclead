package testhelper

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/go-gorp/gorp/v3"
	"github.com/golang-migrate/migrate/v4"
	migratemysql "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

func migrationDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "..", "sql", "migrations")
}

func NewTestDB(t *testing.T) *gorp.DbMap {
	t.Helper()
	ctx := context.Background()

	// 1. MySQLコンテナを起動
	container, err := mysql.Run(ctx,
		"mysql:8.0",
		mysql.WithDatabase("testdb"),
		mysql.WithUsername("root"),
		mysql.WithPassword("password"),
	)
	if err != nil {
		t.Fatalf("mysql container start: %v", err)
	}
	t.Cleanup(func() {
		_ = container.Terminate(ctx)
	})

	// 2. DSNを取得してDBに接続
	dsn, err := container.ConnectionString(ctx, "parseTime=true", "multiStatements=true", "loc=UTC")
	if err != nil {
		t.Fatalf("get connection string: %v", err)
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("sql.Open: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("db.Ping: %v", err)
	}

	// 3. migration を適用
	if err := runMigrations(db); err != nil {
		t.Fatalf("runMigrations: %v", err)
	}

	// 4. gorp.DbMapを返す(本番と同じ設定)
	return &gorp.DbMap{
		Db: db,
		Dialect: gorp.MySQLDialect{
			Engine:   "InnoDB",
			Encoding: "UTF8MB4",
		},
	}
}

func runMigrations(db *sql.DB) error {
	driver, err := migratemysql.WithInstance(db, &migratemysql.Config{})
	if err != nil {
		return fmt.Errorf("migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationDir()),
		"mysql",
		driver,
	)
	if err != nil {
		return fmt.Errorf("migrate new: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate up: %w", err)
	}
	return nil
}
