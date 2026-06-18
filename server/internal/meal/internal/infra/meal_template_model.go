package mealinfra

import (
	"database/sql"
	"time"
)

type MealTemplateModel struct {
	ID            []byte         `db:"id"`
	UserID        []byte         `db:"user_id"`
	Name          string         `db:"name"`
	DisplayOrder  int            `db:"display_order"`
	MealType      string         `db:"meal_type"`
	Calories      int            `db:"calories"`
	ProteinG      sql.NullString `db:"protein_g"`
	FatG          sql.NullString `db:"fat_g"`
	CarbohydrateG sql.NullString `db:"carbohydrate_g"`
	CreatedAt     time.Time      `db:"created_at"`
	UpdatedAt     time.Time      `db:"updated_at"`
}
