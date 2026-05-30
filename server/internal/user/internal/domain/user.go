package userdomain

import (
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

type User struct {
	id           valueobject.UserID
	name         valueobject.String50
	email        valueobject.Email
	passwordHash valueobject.HashedPassword
	birthday     *time.Time
	deletedAt    *time.Time
	createdAt    time.Time
	updatedAt    time.Time
}

func (u *User) ID() valueobject.UserID {
	return u.id
}

func (u *User) Name() valueobject.String50 {
	return u.name
}

func (u *User) Email() valueobject.Email {
	return u.email
}

func (u *User) PasswordHash() valueobject.HashedPassword {
	return u.passwordHash
}

func (u *User) Birthday() *time.Time {
	return u.birthday
}

func (u *User) DeletedAt() *time.Time {
	return u.deletedAt
}

func (u *User) MarkAsDeleted() {
	now := time.Now()
	u.deletedAt = &now
	u.updatedAt = now
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

func CreateUser(
	name valueobject.String50,
	email valueobject.Email,
	passwordHash valueobject.HashedPassword,
	birthday *time.Time,
) *User {
	return &User{
		id:           valueobject.NewPrimaryID[valueobject.UserID](),
		name:         name,
		email:        email,
		passwordHash: passwordHash,
		birthday:     birthday,
		createdAt:    time.Now(),
		updatedAt:    time.Now(),
	}
}

func NewUser(id valueobject.UserID, name valueobject.String50, email valueobject.Email, passwordHash valueobject.HashedPassword, birthday *time.Time, createdAt time.Time, updatedAt time.Time) *User {
	return &User{
		id:           id,
		name:         name,
		email:        email,
		passwordHash: passwordHash,
		birthday:     birthday,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}
