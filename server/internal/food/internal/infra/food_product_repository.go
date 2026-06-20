package foodinfra

import (
	"context"
	"database/sql"

	fooddomain "github.com/Watari995/musclead/internal/food/internal/domain"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/Watari995/musclead/internal/shared/sqlconv"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/go-gorp/gorp/v3"
)

// FoodProductRepository は food_products テーブルの DB 実装。
type FoodProductRepository struct {
	dbmap *gorp.DbMap
}

func NewFoodProductRepository(dbmap *gorp.DbMap) fooddomain.FoodProductRepository {
	return &FoodProductRepository{dbmap: dbmap}
}

func (r *FoodProductRepository) FindByID(ctx context.Context, id valueobject.FoodProductID) (*fooddomain.FoodProduct, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	bytes, err := id.Bytes()
	if err != nil {
		return nil, err
	}
	var row FoodProductModel
	err = q.SelectOne(&row, "SELECT id, barcode, name, calories, protein_g, fat_g, carbohydrate_g, register_source, created_at, updated_at FROM food_products WHERE id = ?", bytes)
	if err != nil {
		return nil, err
	}
	return toFoodProduct(row)
}

func (r *FoodProductRepository) FindAllByBarcode(ctx context.Context, barcode valueobject.Barcode) ([]*fooddomain.FoodProduct, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	var rows []FoodProductModel
	_, err := q.Select(&rows, "SELECT id, barcode, name, calories, protein_g, fat_g, carbohydrate_g, register_source, created_at, updated_at FROM food_products WHERE barcode = ?", barcode.Value())
	if err != nil {
		return nil, err
	}
	result := make([]*fooddomain.FoodProduct, len(rows))
	for i, row := range rows {
		result[i], err = toFoodProduct(row)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (r *FoodProductRepository) FindAllByName(ctx context.Context, name valueobject.String100) ([]*fooddomain.FoodProduct, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	var rows []FoodProductModel
	_, err := q.Select(&rows, "SELECT id, barcode, name, calories, protein_g, fat_g, carbohydrate_g, register_source, created_at, updated_at FROM food_products WHERE name LIKE ?", name.Value()+"%")
	if err != nil {
		return nil, err
	}
	result := make([]*fooddomain.FoodProduct, len(rows))
	for i, row := range rows {
		result[i], err = toFoodProduct(row)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (r *FoodProductRepository) Create(ctx context.Context, foodProduct *fooddomain.FoodProduct) error {
	q := dbtx.Querier(ctx, r.dbmap)
	params, err := buildInsertFoodProductParams(foodProduct)
	if err != nil {
		return err
	}
	_, err = q.Exec("INSERT INTO food_products (id, barcode, name, calories, protein_g, fat_g, carbohydrate_g, register_source, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", params...)
	if err != nil {
		return err
	}
	return nil
}

func buildInsertFoodProductParams(foodProduct *fooddomain.FoodProduct) ([]any, error) {
	bytes, err := foodProduct.ID().Bytes()
	if err != nil {
		return nil, err
	}
	var barcode sql.NullString
	if foodProduct.Barcode() != nil {
		barcode = sqlconv.StringToNullString(foodProduct.Barcode().Value())
	}
	name := foodProduct.Name().Value()
	calories := foodProduct.Calories().Value()
	var proteinG sql.NullString
	if foodProduct.ProteinG() != nil {
		proteinG = sqlconv.DecimalToNullString(foodProduct.ProteinG().Value())
	}
	var fatG sql.NullString
	if foodProduct.FatG() != nil {
		fatG = sqlconv.DecimalToNullString(foodProduct.FatG().Value())
	}
	var carbohydrateG sql.NullString
	if foodProduct.CarbohydrateG() != nil {
		carbohydrateG = sqlconv.DecimalToNullString(foodProduct.CarbohydrateG().Value())
	}
	registerSource := sqlconv.StringToNullString(foodProduct.RegisterSource().Value())
	createdAt := foodProduct.CreatedAt()
	updatedAt := foodProduct.UpdatedAt()
	return []any{bytes, barcode, name, calories, proteinG, fatG, carbohydrateG, registerSource, createdAt, updatedAt}, nil
}

func toFoodProduct(row FoodProductModel) (*fooddomain.FoodProduct, error) {
	id, err := sqlconv.NewPrimaryIDFromBytes[valueobject.FoodProductID](row.ID)
	if err != nil {
		return nil, err
	}
	barcode, err := sqlconv.NewBarcodeFromNullString(row.Barcode)
	if err != nil {
		return nil, err
	}
	name, err := valueobject.NewString100(row.Name)
	if err != nil {
		return nil, err
	}
	calories, err := valueobject.NewNonNegativeInt(row.Calories)
	if err != nil {
		return nil, err
	}
	proteinG, err := sqlconv.NewNonNegativeDecimalFromNullString(row.ProteinG)
	if err != nil {
		return nil, err
	}
	fatG, err := sqlconv.NewNonNegativeDecimalFromNullString(row.FatG)
	if err != nil {
		return nil, err
	}
	carbohydrateG, err := sqlconv.NewNonNegativeDecimalFromNullString(row.CarbohydrateG)
	if err != nil {
		return nil, err
	}
	registerSource, err := valueobject.NewFoodProductRegisterSourceFromString(row.RegisterSource)
	if err != nil {
		return nil, err
	}
	return fooddomain.NewFoodProduct(*id, barcode, *name, *calories, proteinG, fatG, carbohydrateG, *registerSource, row.CreatedAt, row.UpdatedAt), nil
}
