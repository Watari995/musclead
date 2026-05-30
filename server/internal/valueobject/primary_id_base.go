package valueobject

import "github.com/google/uuid"

type PrimaryIDBase struct {
	LiteralBase[string]
}

func newPrimaryIDBase() PrimaryIDBase {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return PrimaryIDBase{
		LiteralBase: LiteralBase[string]{v: id.String()},
	}
}

func newPrimaryIDBaseFromString(s string) (PrimaryIDBase, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return PrimaryIDBase{}, err
	}
	return PrimaryIDBase{
		LiteralBase: LiteralBase[string]{v: id.String()},
	}, nil
}

func (p PrimaryIDBase) Bytes() ([]byte, error) {
	// uuid.Parse() は 16 バイトの配列を返す
	u, err := uuid.Parse(p.v)
	if err != nil {
		return nil, err
	}
	// これを []byte に変換
	return u[:], nil
}
