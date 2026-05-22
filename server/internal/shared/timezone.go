package shared

import "time"

// プロジェクト全体で時刻を UTC で扱う。
// これにより time.Now() などが常に UTC location を返す。
func init() {
	time.Local = time.UTC
}
