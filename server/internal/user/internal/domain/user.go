package userdomain

import (
	"time"

	sharedstorage "github.com/Watari995/musclead/internal/shared/infra/storage"
	"github.com/Watari995/musclead/internal/valueobject"
)

const DefaultProfileImagePath = string(sharedstorage.ImageKindProfile) + "/default.png"

type User struct {
	id               valueobject.UserID
	name             valueobject.String50
	email            valueobject.Email
	passwordHash     valueobject.HashedPassword
	birthday         *time.Time
	profileImagePath string
	deletedAt        *time.Time
	createdAt        time.Time
	updatedAt        time.Time
}

func (u *User) ID() valueobject.UserID {
	return u.id
}

func (u *User) Name() valueobject.String50 {
	return u.name
}

func (u *User) SetName(n valueobject.String50) {
	u.name = n
	u.updatedAt = time.Now()
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

func (u *User) SetBirthday(b *time.Time) {
	u.birthday = b
	u.updatedAt = time.Now()
}

func (u *User) ProfileImagePath() string {
	return u.profileImagePath
}

func (u *User) SetProfileImagePath(p string) {
	u.profileImagePath = p
	u.updatedAt = time.Now()
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

func CreateOnboardingUser(
	name valueobject.String50,
	email valueobject.Email,
	passwordHash valueobject.HashedPassword,
) *User {
	return NewUser(
		valueobject.NewPrimaryID[valueobject.UserID](),
		name,
		email,
		passwordHash,
		nil,
		DefaultProfileImagePath, // 最初はデフォルト画像で作成する
		time.Now(),
		time.Now(),
	)
}

func NewUser(
	id valueobject.UserID,
	name valueobject.String50,
	email valueobject.Email,
	passwordHash valueobject.HashedPassword,
	birthday *time.Time,
	profileImagePath string,
	createdAt time.Time,
	updatedAt time.Time,
) *User {
	return &User{
		id:               id,
		name:             name,
		email:            email,
		passwordHash:     passwordHash,
		birthday:         birthday,
		profileImagePath: profileImagePath,
		createdAt:        createdAt,
		updatedAt:        updatedAt,
	}
}
