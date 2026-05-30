package userinfra

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Watari995/musclead/internal/shared/sqlconv"
	userdomain "github.com/Watari995/musclead/internal/user/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/go-gorp/gorp/v3"
)

type userRepository struct {
	dbmap *gorp.DbMap
}

const upsertUserSQL = `
INSERT INTO users (id, name, email, password_hash, birthday, deleted_at, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    name = VALUES(name),
    email = VALUES(email),
    password_hash = VALUES(password_hash),
    birthday = VALUES(birthday),
    deleted_at = VALUES(deleted_at),
    updated_at = VALUES(updated_at)
`

func (r *userRepository) FindByID(ctx context.Context, id valueobject.UserID) (*userdomain.User, error) {
	bytes, err := id.Bytes()
	if err != nil {
		return nil, err
	}
	var row userModel
	err = r.dbmap.WithContext(ctx).SelectOne(&row,
		"SELECT id, name, email, password_hash, birthday, deleted_at, created_at, updated_at FROM users WHERE id = ? AND deleted_at IS NULL",
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
	var row userModel
	err := r.dbmap.WithContext(ctx).SelectOne(&row,
		"SELECT id, name, email, password_hash, birthday, deleted_at, created_at, updated_at FROM users WHERE email = ? AND deleted_at IS NULL",
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

func (r *userRepository) Save(ctx context.Context, user *userdomain.User) error {
	bytes, err := user.ID().Bytes()
	if err != nil {
		return err
	}
	_, err = r.dbmap.WithContext(ctx).Exec(upsertUserSQL,
		bytes,
		user.Name().Value(),
		user.Email().Value(),
		user.PasswordHash().Value(),
		sqlconv.ToNullTime(user.Birthday()),
		sqlconv.ToNullTime(user.DeletedAt()),
		user.CreatedAt(),
		user.UpdatedAt(),
	)
	return err
}

func toUser(row userModel) (*userdomain.User, error) {
	userIdString, err := sqlconv.UUIDStringFromBytes(row.ID)
	if err != nil {
		return nil, err
	}
	userId, err := valueobject.NewPrimaryIDFromString[valueobject.UserID](userIdString)
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
	return userdomain.NewUser(*userId, *name, *email, *passwordHash, birthday, row.CreatedAt, row.UpdatedAt), nil
}

func NewUserRepository(db *sql.DB) userdomain.UserRepository {
	dbmap := &gorp.DbMap{
		Db:      db,
		Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8MB4"},
	}
	dbmap.AddTableWithName(userModel{}, "users").SetKeys(false, "ID")
	return &userRepository{dbmap: dbmap}
}
