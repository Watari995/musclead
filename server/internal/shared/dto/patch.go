package shareddto

import (
	"bytes"
	"encoding/json"
)

// goでは
// keyがなくてvalueがundefined
// keyがあってvalueがnull
// を区別できないので、それをチェックしてdecodeできるような汎用構造体とメソッドを作成する

type Patch[T any] struct {
	Set   bool
	Null  bool
	Value T
}

func (p *Patch[T]) UnmarshalJSON(data []byte) error {
	// ここに入った段階でKeyはある
	p.Set = true

	// nullが入っていたらNullをtrueにする
	if bytes.Equal(data, []byte("null")) {
		p.Null = true
		return nil
	}
	return json.Unmarshal(data, &p.Value)
}
