package userinfra

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/Watari995/musclead/internal/shared/sqlconv"
	userdomain "github.com/Watari995/musclead/internal/user/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/go-gorp/gorp/v3"
)

type userRepository struct {
	dbmap *gorp.DbMap
}

const upsertUserSQL = `
INSERT INTO users (id, name, email, password_hash, birthday, profile_image_path, deleted_at, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    name = VALUES(name),
    email = VALUES(email),
    password_hash = VALUES(password_hash),
    birthday = VALUES(birthday),
    profile_image_path = VALUES(profile_image_path),
    deleted_at = VALUES(deleted_at),
    updated_at = VALUES(updated_at)
`

func (r *userRepository) FindByID(ctx context.Context, id valueobject.UserID) (*userdomain.User, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	bytes, err := id.Bytes()
	if err != nil {
		return nil, err
	}
	var row UserModel
	err = q.SelectOne(&row,
		"SELECT id, name, email, password_hash, birthday, profile_image_path, deleted_at, created_at, updated_at FROM users WHERE id = ? AND deleted_at IS NULL",
		bytes,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toUser(row)
}

func (r *userRepository) FindByEmail(ctx context.Context, email valueobject.Email) (*userdomain.User, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	var row UserModel
	err := q.SelectOne(&row,
		"SELECT id, name, email, password_hash, birthday, profile_image_path, deleted_at, created_at, updated_at FROM users WHERE email = ? AND deleted_at IS NULL",
		email.Value(),
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toUser(row)
}

type getAllUserIDsRow struct {
	ID []byte `db:"id"`
}

func (r *userRepository) GetAllUserIDs(ctx context.Context) ([]valueobject.UserID, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	var rows []getAllUserIDsRow
	if _, err := q.Select(&rows, "SELECT id FROM users WHERE deleted_at IS NULL"); err != nil {
		return nil, err
	}
	result := make([]valueobject.UserID, 0, len(rows))
	for _, row := range rows {
		userID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.UserID](row.ID)
		if err != nil {
			return nil, err
		}
		result = append(result, *userID)
	}
	return result, nil
}

func (r *userRepository) Save(ctx context.Context, user *userdomain.User) error {
	q := dbtx.Querier(ctx, r.dbmap)
	bytes, err := user.ID().Bytes()
	if err != nil {
		return err
	}
	_, err = q.Exec(upsertUserSQL,
		bytes,
		user.Name().Value(),
		user.Email().Value(),
		user.PasswordHash().Value(),
		sqlconv.ToNullTime(user.Birthday()),
		user.ProfileImagePath(),
		sqlconv.ToNullTime(user.DeletedAt()),
		user.CreatedAt(),
		user.UpdatedAt(),
	)
	return err
}

func toUser(row UserModel) (*userdomain.User, error) {
	userID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.UserID](row.ID)
	if err != nil {
		return nil, err
	}
	name, err := valueobject.NewString50(row.Name)
	if err != nil {
		return nil, err
	}
	email, err := valueobject.NewEmail(row.Email)
	if err != nil {
		return nil, err
	}
	passwordHash, err := valueobject.NewHashedPassword(row.PasswordHash)
	if err != nil {
		return nil, err
	}
	birthday := sqlconv.FromNullTime(row.Birthday)
	return userdomain.NewUser(*userID, *name, *email, *passwordHash, birthday, row.ProfileImagePath, row.CreatedAt, row.UpdatedAt), nil
}

func NewUserRepository(dbmap *gorp.DbMap) userdomain.UserRepository {
	return &userRepository{dbmap: dbmap}
}
