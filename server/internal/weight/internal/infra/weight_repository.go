package weightinfra

import (
	"context"

	weightdomain "github.com/Watari995/musclead/internal/weight/internal/domain"
	"github.com/go-gorp/gorp/v3"
)

type weightRepository struct {
	dbmap *gorp.DbMap
}

func NewWeightRepository(dbmap *gorp.DbMap) weightdomain.WeightRepository {
	return &weightRepository{dbmap: dbmap}
}

const upsertWeightSQL = `
INSERT INTO weights (id, user_id, weight_kg, body_fat_percentage, skeletal_muscle_kg, measured_at, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    weight_kg = VALUES(weight_kg),
    body_fat_percentage = VALUES(body_fat_percentage),
    skeletal_muscle_kg = VALUES(skeletal_muscle_kg),
    measured_at = VALUES(measured_at),
    updated_at = VALUES(updated_at)
`

func (r *weightRepository) FindByIDAndUserID(ctx context.Context, id valueobject.WeightID, userID valueobject.UserID) (*weightdomain.Weight, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	idBytes, err := id.Bytes()
	if err != nil {
		return nil, err
	}
	userIDBytes, err := userID.Bytes()
	if err != nil {
		return nil, err
	}
}
