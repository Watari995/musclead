package valueobject

import "github.com/google/uuid"

type PrimaryIdBase struct {
	LiteralBase[string]
}

func newPrimaryIdBase() PrimaryIdBase {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return PrimaryIdBase{
		LiteralBase: LiteralBase[string]{v: id.String()},
	}
}

func newPrimaryIdBaseFromString(s string) (PrimaryIdBase, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return PrimaryIdBase{}, err
	}
	return PrimaryIdBase{
		LiteralBase: LiteralBase[string]{v: id.String()},
	}, nil
}
