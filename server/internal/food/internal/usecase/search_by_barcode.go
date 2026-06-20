package foodusecase

// SearchByBarcode は自社 DB → Open Food Facts の順でバーコード検索する。
// 自社 DB になければ外部 API を呼び、ヒットした場合は DB にキャッシュして返す。
// どちらもなければ not found エラーを返す。
type SearchByBarcode struct {
}

func NewSearchByBarcode() *SearchByBarcode {
	panic("not implemented")
}
