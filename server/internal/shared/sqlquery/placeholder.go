// Package sqlquery は raw SQL を組み立てる際に必要となる、 文法的なヘルパーを集める。
//
// gorp のように構造体マッピングまで担うラッパーとは違って、
// 「SQL 構文を作るためだけの薄い文字列操作」 をここに置く。
// 値オブジェクトの詰め替えは sqlconv、 TX 制御は dbtx、 SQL 文の組み立て補助はここ、 という分業。
package sqlquery

import "strings"

// InPlaceholders は IN 句に複数値を流し込みたい時の placeholder 文字列と
// args スライスを同時に作る。
//
// 例: items が 3要素なら "?,?,?" と [a, b, c] を返す。
// 呼び出し側は: q.Select(&rows, "... WHERE id IN ("+ph+")", args...)
//
// 想定用途: 階層的な集約をまとめ読みする時の「親 ID リストで子をまとめて取る」 クエリ。
// (N+1 を避けるため repository 内で頻出。)
//
// items が空の場合は空文字列と nil を返すので、 呼び出し側は事前に len チェックして
// クエリ自体を発行しないこと。 IN () は MySQL では構文エラーになる。
func InPlaceholders(items [][]byte) (string, []any) {
	if len(items) == 0 {
		return "", nil
	}
	// "?,?,?,?,..." を作って末尾のカンマを落とす
	placeholders := strings.TrimRight(strings.Repeat("?,", len(items)), ",")
	args := make([]any, 0, len(items))
	for _, item := range items {
		args = append(args, item)
	}
	return placeholders, args
}
