package sqlconv

import (
	"database/sql"
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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

func DecimalFromNullString(s sql.NullString) (*decimal.Decimal, error) {
	if !s.Valid {
		return nil, nil
	}
	d, err := decimal.NewFromString(s.String)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func DecimalToNullString(d decimal.Decimal) sql.NullString {
	return sql.NullString{
		String: d.String(),
		Valid:  true,
	}
}

func NewNonNegativeDecimalFromNullString(s sql.NullString) (*valueobject.NonNegativeDecimal, error) {
	d, err := DecimalFromNullString(s)
	if err != nil {
		return nil, err
	}
	if d != nil {
		vo, err := valueobject.NewNonNegativeDecimal(*d)
		if err != nil {
			return nil, err
		}
		return vo, nil
	}
	return nil, nil
}
