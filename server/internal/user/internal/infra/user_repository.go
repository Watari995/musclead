package userinfra

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Watari995/musclead/internal/shared/sqlconv"
	userdomain "github.com/Watari995/musclead/internal/user/internal/domain"
	userdbgen "github.com/Watari995/musclead/internal/user/internal/infra/dbgen"
	"github.com/Watari995/musclead/internal/valueobject"
)

type userRepository struct {
	db *userdbgen.Queries
}

func (r *userRepository) FindByID(ctx context.Context, id valueobject.UserID) (*userdomain.User, error) {
	bytes, err := id.Bytes()
	if err != nil {
		return nil, err
	}
	userRow, err := r.db.FindUserByID(ctx, bytes)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toUser(userRow)
}

func (r *userRepository) FindByEmail(ctx context.Context, email valueobject.Email) (*userdomain.User, error) {
	userRow, err := r.db.FindUserByEmail(ctx, email.Value())
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toUser(userRow)
}

func (r *userRepository) Save(ctx context.Context, user *userdomain.User) error {
	params, err := toUpsertUserParams(user)
	if err != nil {
		return err
	}
	return r.db.UpsertUser(ctx, params)
}

// toUser is a helper function to convert a userdbgen.User to a userdomain.User
func toUser(userRow userdbgen.User) (*userdomain.User, error) {
	userIdString, err := sqlconv.UUIDStringFromBytes(userRow.ID)
	if err != nil {
		return nil, err
	}
	userId, err := valueobject.NewPrimaryIdFromString[valueobject.UserID](userIdString)
	if err != nil {
		return nil, err
	}
	name, err := valueobject.NewString50(userRow.Name)
	if err != nil {
		return nil, err
	}
	email, err := valueobject.NewEmail(userRow.Email)
	if err != nil {
		return nil, err
	}
	passwordHash, err := valueobject.NewHashedPassword(userRow.PasswordHash)
	if err != nil {
		return nil, err
	}
	birthday := sqlconv.FromNullTime(userRow.Birthday)
	return userdomain.NewUser(*userId, *name, *email, *passwordHash, birthday, userRow.CreatedAt, userRow.UpdatedAt), nil
}

func toUpsertUserParams(user *userdomain.User) (userdbgen.UpsertUserParams, error) {
	bytes, err := user.ID().Bytes()
	if err != nil {
		return userdbgen.UpsertUserParams{}, err
	}
	return userdbgen.UpsertUserParams{
		ID:           bytes,
		Name:         user.Name().Value(),
		Email:        user.Email().Value(),
		PasswordHash: user.PasswordHash().Value(),
		Birthday:     sqlconv.ToNullTime(user.Birthday()),
		DeletedAt:    sqlconv.ToNullTime(user.DeletedAt()),
		CreatedAt:    user.CreatedAt(),
		UpdatedAt:    user.UpdatedAt(),
	}, nil
}

func NewUserRepository(db *userdbgen.Queries) userdomain.UserRepository {
	return &userRepository{db: db}
}
