package mealinfra

import (
	"database/sql"
	"time"
)

type MealModel struct {
	ID            []byte         `db:"id"`
	UserID        []byte         `db:"user_id"`
	EatenAt       time.Time      `db:"eaten_at"`
	MealType      string         `db:"meal_type"`
	Calories      int          `db:"calories"`
	ProteinG      sql.NullString `db:"protein_g"`
	FatG          sql.NullString `db:"fat_g"`
	CarbohydrateG sql.NullString `db:"carbohydrate_g"`
	Memo          sql.NullString `db:"memo"`
	CreatedAt     time.Time      `db:"created_at"`
	UpdatedAt     time.Time      `db:"updated_at"`
}

type MealPhotoModel struct {
	ID           []byte    `db:"id"`
	MealID       []byte    `db:"meal_id"`
	ImagePath    string    `db:"image_path"`
	DisplayOrder int     `db:"display_order"`
	CreatedAt    time.Time `db:"created_at"`
}
