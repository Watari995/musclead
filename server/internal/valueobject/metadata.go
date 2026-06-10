package valueobject

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Metadata は JSON column を Go の map として扱う VO。
// driver.Valuer / sql.Scanner を実装することで、 gorp / database/sql が
// DB との読み書き時に自動でシリアライズ / デシリアライズ する。
type Metadata map[string]any

// Value は driver.Valuer interface の実装。
// gorp が DB 書き込み時に呼び出し、 map を JSON []byte に変換する。
// nil の場合は DB NULL として扱う。
func (m Metadata) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

// Scan は sql.Scanner interface の実装。
// gorp が DB 読み込み時に呼び出し、 JSON []byte を map に復元する。
// DB NULL の場合は nil map にする。
func (m *Metadata) Scan(value any) error {
	if value == nil {
		*m = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte, got %T", value)
	}
	return json.Unmarshal(bytes, m)
}
