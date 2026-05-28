package storage

import (
	"fmt"
	"strings"
)

func BuildImageURL(baseURL string, path string) string {
	// 先頭のスラッシュを削除する(あれば)
	path = strings.TrimPrefix(path, "/")
	// baseURL の末尾にスラッシュがなければ追加する
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	return fmt.Sprintf("%s%s", baseURL, path)
}
