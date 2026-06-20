// Package sqlconv は DB の sql.NullXxx / 生バイト型と、 ドメイン層の VO / time.Time / decimal などを相互変換するヘルパー集。
//
// repository 実装が NULL 列や UUID バイト列を扱うたびに同じ詰め替えを書かないよう、
// ここに一箇所で寄せる。 「列 = NULL 許容」 なら *valueobject.X / *time.Time のようなポインタ表現、
// 「列 = NOT NULL」 なら値型をそのまま使う、 という Go 慣習に従う。
//
// 命名規約:
//   - To<NullType>           : VO/値 → sql.NullXxx (UPDATE/INSERT 用)
//   - From<NullType>         : sql.NullXxx → VO/値 (SELECT 用)
//   - New<VO>FromNullXxx     : NULL 起点で VO を生成(検証込み、 error を返す)
//
// 新しい VO を NULL 列で扱いたい時はこのファイルに対応ペアを追加すること。
package sqlconv

import (
	"database/sql"
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ─── time ────────────────────────────────────────────────

func ToNullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

func FromNullTime(t sql.NullTime) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

// ─── string / String1000 ─────────────────────────────────

func StringToNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: true}
}

func StringPtrToNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}

func NewStringFromNullString(s sql.NullString) *string {
	if !s.Valid {
		return nil
	}
	return &s.String
}

// String1000ToNullString は *String1000 を sql.NullString に詰める。
// nil の時は Valid:false(= DB NULL) を返す。
func String1000ToNullString(v *valueobject.String1000) sql.NullString {
	if v == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: v.Value(), Valid: true}
}

// NewString1000FromNullString は sql.NullString → *String1000 への復元。
// NULL なら nil を返す。 値があるが VO 制約に違反する場合は error。
func NewString1000FromNullString(s sql.NullString) (*valueobject.String1000, error) {
	if !s.Valid {
		return nil, nil
	}
	return valueobject.NewString1000(s.String)
}

// ─── barcode / Barcode ───────────────────────────────────

func NewBarcodeFromNullString(s sql.NullString) (*valueobject.Barcode, error) {
	if !s.Valid {
		return nil, nil
	}
	return valueobject.NewBarcode(s.String)
}

// ─── decimal / NonNegativeDecimal ───────────────────────

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
	return sql.NullString{String: d.String(), Valid: true}
}

func NewNonNegativeDecimalFromNullString(s sql.NullString) (*valueobject.NonNegativeDecimal, error) {
	d, err := DecimalFromNullString(s)
	if err != nil {
		return nil, err
	}
	if d == nil {
		return nil, nil
	}
	return valueobject.NewNonNegativeDecimal(*d)
}

// ─── int / NonNegativeInt ───────────────────────────────

// NonNegativeIntToNullInt32 は *NonNegativeInt を sql.NullInt32 に詰める。
// nil の時は Valid:false(= DB NULL)。
func NonNegativeIntToNullInt32(v *valueobject.NonNegativeInt) sql.NullInt32 {
	if v == nil {
		return sql.NullInt32{}
	}
	return sql.NullInt32{Int32: int32(v.Value()), Valid: true}
}

// NewNonNegativeIntFromNullInt32 は sql.NullInt32 → *NonNegativeInt への復元。
// NULL なら nil。 値が範囲外なら error。
func NewNonNegativeIntFromNullInt32(v sql.NullInt32) (*valueobject.NonNegativeInt, error) {
	if !v.Valid {
		return nil, nil
	}
	return valueobject.NewNonNegativeInt(int(v.Int32))
}

// ─── percentage / Percentage ─────────────────────────────

func NewPercentageFromNullString(s sql.NullString) (*valueobject.Percentage, error) {
	if !s.Valid {
		return nil, nil
	}
	d, err := decimal.NewFromString(s.String)
	if err != nil {
		return nil, err
	}
	return valueobject.NewPercentage(d)
}

func PercentageToNullString(p valueobject.Percentage) sql.NullString {
	return sql.NullString{String: p.Value().String(), Valid: true}
}

// ─── weight_kg / WeightKg ───────────────────────────────

func NewWeightKgFromString(s string) (*valueobject.WeightKg, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return nil, err
	}
	return valueobject.NewWeightKg(d)
}

func WeightKgToNullString(w valueobject.WeightKg) sql.NullString {
	return sql.NullString{String: w.Value().String(), Valid: true}
}

func NewWeightKgFromNullString(s sql.NullString) (*valueobject.WeightKg, error) {
	if !s.Valid {
		return nil, nil
	}
	return NewWeightKgFromString(s.String)
}

// ─── UUID ───────────────────────────────────────────────

func UUIDStringFromBytes(b []byte) (string, error) {
	u, err := uuid.FromBytes(b)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

// byte->PrimaryIDVOに変換する
func NewPrimaryIDFromBytes[T valueobject.PrimaryID](b []byte) (*T, error) {
	s, err := UUIDStringFromBytes(b)
	if err != nil {
		return nil, err
	}
	return valueobject.NewPrimaryIDFromString[T](s)
}

func NewPrimaryIDFromNullableBytes[T valueobject.PrimaryID](b []byte) (*T, error) {
	// スライスは初期値がnil
	if b == nil {
		return nil, nil
	}
	return NewPrimaryIDFromBytes[T](b)
}

// UUID文字列 -> 16 bytes binary
func UUIDStringToBytes(s string) ([]byte, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return nil, err
	}
	return u[:], nil
}

// nullable uuid -> []byte
func NewBytesFromNullablePrimaryID[T valueobject.PrimaryID](id *T) ([]byte, error) {
	if id == nil {
		return nil, nil
	}
	bytes, err := (*id).Bytes()
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// ---- url / URL ─────────────────────────────────────────────

func NewURLFromNullString(s sql.NullString) (*valueobject.URL, error) {
	if !s.Valid {
		return nil, nil
	}
	return valueobject.NewURL(s.String)
}

func URLPtrToNullString(u *valueobject.URL) sql.NullString {
	if u == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: u.Value(), Valid: true}
}
