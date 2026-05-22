package valueobject

type PrimaryId interface {
	// PrimaryIdBase を埋め込んだ型を指定
	~struct{ PrimaryIdBase }
}

func NewPrimaryId[T PrimaryId]() T {
	return T{PrimaryIdBase: newPrimaryIdBase()}
}

func NewPrimaryIdFromString[T PrimaryId](s string) (*T, error) {
	p, err := newPrimaryIdBaseFromString(s)
	if err != nil {
		return nil, err
	}
	return &T{PrimaryIdBase: p}, nil
}
