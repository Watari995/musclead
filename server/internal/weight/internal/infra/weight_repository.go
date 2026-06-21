package weightinfra

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Watari995/musclead/internal/pagination"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/Watari995/musclead/internal/shared/sqlconv"
	"github.com/Watari995/musclead/internal/valueobject"
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
	var row WeightModel
	err = q.SelectOne(&row,
		"SELECT id, user_id, weight_kg, body_fat_percentage, skeletal_muscle_kg, measured_at, created_at, updated_at FROM weights WHERE id = ? AND user_id = ?",
		idBytes, userIDBytes,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toEntity(row)
}

func (r *weightRepository) FindAllByUserIDWithOffsetPagination(ctx context.Context, userID valueobject.UserID, limit int, offset int) ([]*weightdomain.Weight, pagination.OffsetPaginator, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	userIDBytes, err := userID.Bytes()
	if err != nil {
		return nil, pagination.OffsetPaginator{}, err
	}
	var rows []WeightModel
	_, err = q.Select(&rows,
		"SELECT id, user_id, weight_kg, body_fat_percentage, skeletal_muscle_kg, measured_at, created_at, updated_at FROM weights WHERE user_id = ? ORDER BY measured_at DESC LIMIT ? OFFSET ?",
		userIDBytes, int32(limit), int32(offset),
	)
	if err != nil {
		return nil, pagination.OffsetPaginator{}, err
	}
	total, err := q.SelectInt(
		"SELECT COUNT(*) FROM weights WHERE user_id = ?", userIDBytes,
	)
	if err != nil {
		return nil, pagination.OffsetPaginator{}, err
	}
	paginator := pagination.NewOffsetPaginator(int(total), offset, limit)
	if len(rows) == 0 {
		return []*weightdomain.Weight{}, paginator, nil
	}
	result := make([]*weightdomain.Weight, len(rows))
	for i, row := range rows {
		weight, err := toEntity(row)
		if err != nil {
			return nil, pagination.OffsetPaginator{}, err
		}
		result[i] = weight
	}
	return result, paginator, nil
}

func (r *weightRepository) ExistsByUserIDAndMeasuredAt(ctx context.Context, userID valueobject.UserID, measuredAt time.Time) (bool, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	userIDBytes, err := userID.Bytes()
	if err != nil {
		return false, err
	}
	var exists bool
	err = q.SelectOne(&exists, "SELECT EXISTS(SELECT 1 FROM weights WHERE user_id = ? AND measured_at = ?)", userIDBytes, measuredAt)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *weightRepository) FindAllByUserIDAndPeriod(ctx context.Context, userID valueobject.UserID, from, to time.Time) ([]*weightdomain.Weight, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	userIDBytes, err := userID.Bytes()
	if err != nil {
		return nil, err
	}
	var rows []WeightModel
	_, err = q.Select(&rows, "SELECT id, user_id, weight_kg, body_fat_percentage, skeletal_muscle_kg, measured_at, created_at, updated_at FROM weights WHERE user_id = ? AND measured_at BETWEEN ? AND ? ORDER BY measured_at ASC", userIDBytes, from, to)
	if err != nil {
		return nil, err
	}
	result := make([]*weightdomain.Weight, len(rows))
	for i, row := range rows {
		weight, err := toEntity(row)
		if err != nil {
			return nil, err
		}
		result[i] = weight
	}
	return result, nil
}

func (r *weightRepository) Save(ctx context.Context, weight *weightdomain.Weight) error {
	q := dbtx.Querier(ctx, r.dbmap)
	params, err := buildUpsertWeightParams(weight)
	if err != nil {
		return err
	}
	if _, err := q.Exec(upsertWeightSQL, params...); err != nil {
		return err
	}
	return nil
}

func (r *weightRepository) DeleteByID(ctx context.Context, id valueobject.WeightID) error {
	q := dbtx.Querier(ctx, r.dbmap)
	bytes, err := id.Bytes()
	if err != nil {
		return err
	}
	if _, err := q.Exec("DELETE FROM weights WHERE id = ?", bytes); err != nil {
		return err
	}
	return nil
}

func toEntity(row WeightModel) (*weightdomain.Weight, error) {
	weightID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.WeightID](row.ID)
	if err != nil {
		return nil, err
	}
	userID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.UserID](row.UserID)
	if err != nil {
		return nil, err
	}
	weightKg, err := sqlconv.NewWeightKgFromString(row.WeightKg)
	if err != nil {
		return nil, err
	}
	bodyFatPercentage, err := sqlconv.NewPercentageFromNullString(row.BodyFatPercentage)
	if err != nil {
		return nil, err
	}
	skeletalMuscleKg, err := sqlconv.NewWeightKgFromNullString(row.SkeletalMuscleKg)
	if err != nil {
		return nil, err
	}
	return weightdomain.NewWeight(
		*weightID,
		*userID,
		*weightKg,
		bodyFatPercentage,
		skeletalMuscleKg,
		row.MeasuredAt,
		row.CreatedAt,
		row.UpdatedAt,
	), nil
}

func buildUpsertWeightParams(weight *weightdomain.Weight) ([]any, error) {
	idBytes, err := weight.ID().Bytes()
	if err != nil {
		return nil, err
	}
	userIDBytes, err := weight.UserID().Bytes()
	if err != nil {
		return nil, err
	}
	var bodyFatPercentage sql.NullString
	if weight.BodyFatPercentage() != nil {
		bodyFatPercentage = sqlconv.PercentageToNullString(*weight.BodyFatPercentage())
	}
	var skeletalMuscleKg sql.NullString
	if weight.SkeletalMuscleKg() != nil {
		skeletalMuscleKg = sqlconv.WeightKgToNullString(*weight.SkeletalMuscleKg())
	}
	return []any{
		idBytes,
		userIDBytes,
		weight.WeightKg().Value(),
		bodyFatPercentage,
		skeletalMuscleKg,
		weight.MeasuredAt(),
		weight.CreatedAt(),
		weight.UpdatedAt(),
	}, nil
}
