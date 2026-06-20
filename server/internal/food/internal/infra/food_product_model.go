package foodinfra

import (
	"database/sql"
	"time"
)

// FoodProductModel は food_products テーブルの gorp マッピングモデル。
type FoodProductModel struct {
	ID             []byte         `db:"id"`
	Barcode        sql.NullString `db:"barcode"`
	Name           string         `db:"name"`
	Calories       int            `db:"calories"`
	ProteinG       sql.NullString `db:"protein_g"`
	FatG           sql.NullString `db:"fat_g"`
	CarbohydrateG  sql.NullString `db:"carbohydrate_g"`
	RegisterSource string         `db:"register_source"`
	CreatedAt      time.Time      `db:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at"`
}
