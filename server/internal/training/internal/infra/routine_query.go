package traininginfra

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Watari995/musclead/internal/pagination"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/Watari995/musclead/internal/shared/sqlconv"
	"github.com/Watari995/musclead/internal/shared/sqlquery"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/go-gorp/gorp/v3"
	"github.com/samber/lo"
)

type routineQueryService struct {
	dbmap *gorp.DbMap
}

func NewRoutineQueryService(dbmap *gorp.DbMap) trainingdomain.RoutineQueryService {
	return &routineQueryService{dbmap: dbmap}
}

const findRoutineByIDAndUserIDSQL = `
SELECT
	r.id, r.user_id, r.name, r.created_at, r.updated_at,
	re.id AS re_id,
	re.exercise_id AS re_exercise_id,
	e.name AS re_exercise_name,
	re.display_order AS re_display_order
FROM routines r
LEFT JOIN routine_exercises re ON re.routine_id = r.id
LEFT JOIN exercises e ON e.id = re.exercise_id
WHERE r.id = ? AND r.user_id = ?
ORDER BY re.display_order ASC
`

const findRoutineListSQL = `
SELECT id, user_id, name, created_at, updated_at
FROM routines
WHERE user_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?
`

const countRoutinesByUserIDSQL = `SELECT COUNT(*) FROM routines WHERE user_id = ?`

func buildLoadRoutineExercisesByIDsSQL(routineIDs [][]byte) (string, []any) {
	placeholders, args := sqlquery.InPlaceholders(routineIDs)
	return fmt.Sprintf(`
	SELECT
		re.id AS re_id,
		re.routine_id,
		re.exercise_id,
		e.name AS exercise_name,
		re.display_order
	FROM routine_exercises re
	INNER JOIN exercises e ON e.id = re.exercise_id
	WHERE re.routine_id IN (%s)
	ORDER BY re.routine_id, re.display_order ASC
	`, placeholders), args
}

// routineDetailRow は findRoutineByIDAndUserIDSQL の 1 行を受ける。
// re_* / re_exercise_name / re_display_order は LEFT JOIN で NULL になりうるため
// nullable で受ける(BINARY 列は []byte の長さ 0 で NULL 判定)。
type routineDetailRow struct {
	RoutineID    []byte         `db:"id"`
	UserID       []byte         `db:"user_id"`
	Name         string         `db:"name"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
	REID         []byte         `db:"re_id"`
	REExerciseID []byte         `db:"re_exercise_id"`
	ExerciseName sql.NullString `db:"re_exercise_name"`
	DisplayOrder sql.NullInt32  `db:"re_display_order"`
}

func (r *routineQueryService) FindByIDAndUserID(ctx context.Context, id valueobject.RoutineID, userID valueobject.UserID) (*trainingdomain.RoutineView, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	idBytes, err := id.Bytes()
	if err != nil {
		return nil, err
	}
	userIDBytes, err := userID.Bytes()
	if err != nil {
		return nil, err
	}
	var rows []routineDetailRow
	if _, err = q.Select(&rows, findRoutineByIDAndUserIDSQL, idBytes, userIDBytes); err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}
	// 親を取得する (LEFT JOIN で routineに一つに対して複数のレコードが返ってくるため、routineは最初のみでいい)
	first := rows[0]
	routineID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.RoutineID](first.RoutineID)
	if err != nil {
		return nil, err
	}
	routineUserID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.UserID](first.UserID)
	if err != nil {
		return nil, err
	}
	routineName, err := valueobject.NewString50(first.Name)
	if err != nil {
		return nil, err
	}
	exerciseViews := make([]trainingdomain.RoutineExerciseView, 0, len(rows))
	for _, row := range rows {
		if len(row.REID) == 0 {
			continue
		}
		exerciseView, err := toRoutineExerciseViewFromDetail(row)
		if err != nil {
			return nil, err
		}
		exerciseViews = append(exerciseViews, *exerciseView)
	}
	return &trainingdomain.RoutineView{
		ID:               *routineID,
		UserID:           *routineUserID,
		Name:             *routineName,
		CreatedAt:        first.CreatedAt,
		UpdatedAt:        first.UpdatedAt,
		RoutineExercises: exerciseViews,
	}, nil
}

func toRoutineExerciseViewFromDetail(row routineDetailRow) (*trainingdomain.RoutineExerciseView, error) {
	reID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.RoutineExerciseID](row.REID)
	if err != nil {
		return nil, err
	}
	exID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.ExerciseID](row.REExerciseID)
	if err != nil {
		return nil, err
	}
	exName, err := valueobject.NewString50(row.ExerciseName.String)
	if err != nil {
		return nil, err
	}
	displayOrder, err := sqlconv.NewNonNegativeIntFromNullInt32(row.DisplayOrder)
	if err != nil {
		return nil, err
	}
	return &trainingdomain.RoutineExerciseView{
		ID:           *reID,
		ExerciseID:   *exID,
		ExerciseName: *exName,
		DisplayOrder: *displayOrder,
	}, nil
}

// ============== LIST QUERY ==============
type routineListRow struct {
	ID        []byte    `db:"id"`
	UserID    []byte    `db:"user_id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type routineExerciseListRow struct {
	ID           []byte `db:"re_id"`
	RoutineID    []byte `db:"routine_id"`
	ExerciseID   []byte `db:"exercise_id"`
	ExerciseName string `db:"exercise_name"`
	DisplayOrder int32  `db:"display_order"`
}

func (r *routineQueryService) FindAllByUserIDWithOffsetPagination(ctx context.Context, userID valueobject.UserID, limit, offset int) ([]*trainingdomain.RoutineView, pagination.OffsetPaginator, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	userIDBytes, err := userID.Bytes()
	if err != nil {
		return nil, pagination.OffsetPaginator{}, err
	}
	// 1. routineをlimitとoffsetで取得する
	var rows []routineListRow
	if _, err = q.Select(&rows, findRoutineListSQL, userIDBytes, limit, offset); err != nil {
		return nil, pagination.OffsetPaginator{}, err
	}
	total, err := q.SelectInt(countRoutinesByUserIDSQL, userIDBytes)
	if err != nil {
		return nil, pagination.OffsetPaginator{}, err
	}
	paginator := pagination.NewOffsetPaginator(int(total), offset, limit)
	if len(rows) == 0 {
		return []*trainingdomain.RoutineView{}, paginator, nil
	}
	routineIDs := lo.Map(rows, func(row routineListRow, _ int) []byte {
		return row.ID
	})
	// 2. routineに紐づいたroutine_exercisesとexercisesを一括取得する
	var exerciseRows []routineExerciseListRow
	sqlStr, args := buildLoadRoutineExercisesByIDsSQL(routineIDs)
	if _, err = q.Select(&exerciseRows, sqlStr, args...); err != nil {
		return nil, pagination.OffsetPaginator{}, err
	}
	// 3. 該当のroutineのexerciseを取得する
	exerciseMap := make(map[string][]routineExerciseListRow)
	for _, row := range exerciseRows {
		routineID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.RoutineID](row.RoutineID)
		if err != nil {
			return nil, pagination.OffsetPaginator{}, err
		}
		key := routineID.Value()
		exerciseMap[key] = append(exerciseMap[key], row)
	}
	routineViews := make([]*trainingdomain.RoutineView, 0, len(rows))
	for _, row := range rows {
		routineView, err := toRoutineViewFromListRow(row, exerciseMap)
		if err != nil {
			return nil, pagination.OffsetPaginator{}, err
		}
		routineViews = append(routineViews, routineView)
	}
	return routineViews, paginator, nil
}

func toRoutineViewFromListRow(row routineListRow, exerciseMap map[string][]routineExerciseListRow) (*trainingdomain.RoutineView, error) {
	routineID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.RoutineID](row.ID)
	if err != nil {
		return nil, err
	}
	routineUserID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.UserID](row.UserID)
	if err != nil {
		return nil, err
	}
	name, err := valueobject.NewString50(row.Name)
	if err != nil {
		return nil, err
	}
	exerciseViews := make([]trainingdomain.RoutineExerciseView, 0, len(exerciseMap[routineID.Value()]))
	for _, exerciseRow := range exerciseMap[routineID.Value()] {
		exerciseView, err := toRoutineExerciseViewFromListRow(exerciseRow)
		if err != nil {
			return nil, err
		}
		exerciseViews = append(exerciseViews, *exerciseView)
	}
	return &trainingdomain.RoutineView{
		ID:               *routineID,
		UserID:           *routineUserID,
		Name:             *name,
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
		RoutineExercises: exerciseViews,
	}, nil
}

func toRoutineExerciseViewFromListRow(row routineExerciseListRow) (*trainingdomain.RoutineExerciseView, error) {
	reID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.RoutineExerciseID](row.ID)
	if err != nil {
		return nil, err
	}
	exerciseID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.ExerciseID](row.ExerciseID)
	if err != nil {
		return nil, err
	}
	exerciseName, err := valueobject.NewString50(row.ExerciseName)
	if err != nil {
		return nil, err
	}
	displayOrder, err := valueobject.NewNonNegativeInt(int(row.DisplayOrder))
	if err != nil {
		return nil, err
	}
	return &trainingdomain.RoutineExerciseView{
		ID:           *reID,
		ExerciseID:   *exerciseID,
		ExerciseName: *exerciseName,
		DisplayOrder: *displayOrder,
	}, nil
}
