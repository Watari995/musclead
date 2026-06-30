package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"time"

	"github.com/go-gorp/gorp/v3"
	_ "github.com/go-sql-driver/mysql"

	"github.com/Watari995/musclead/internal/meal"
	"github.com/Watari995/musclead/internal/notification"
	"github.com/Watari995/musclead/internal/training"
	"github.com/Watari995/musclead/internal/user"
	"github.com/Watari995/musclead/internal/weight"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := newDB()
	if err != nil {
		slog.Error("failed to connect to db", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	dbmap := &gorp.DbMap{
		Db:      db,
		Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8MB4"},
	}

	userModule := user.NewModule(dbmap, nil, nil)
	trainingModule := training.NewModule(dbmap, nil, nil)
	mealModule := meal.NewModule(dbmap, nil, nil)
	weightModule := weight.NewModule(dbmap, nil)
	notifModule := notification.NewModule(dbmap)

	run(ctx,
		userModule.UserQuery(),
		trainingModule.TrainingQuery(),
		mealModule.MealQuery(),
		weightModule.WeightQuery(),
		notifModule.NotificationCommand(),
	)
}

func newDB() (*sql.DB, error) {
	dsn := os.Getenv("DB_DSN")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
