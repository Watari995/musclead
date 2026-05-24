package sqlconv

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

func ToNullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{
		Time:  *t,
		Valid: true,
	}
}

func FromNullTime(t sql.NullTime) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

func UUIDStringFromBytes(b []byte) (string, error) {
	u, err := uuid.FromBytes(b)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
