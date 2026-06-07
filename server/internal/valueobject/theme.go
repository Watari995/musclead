package valueobject

import "errors"

type ThemeType string

const (
	ThemeLight  ThemeType = "light"
	ThemeDark   ThemeType = "dark"
	ThemeSystem ThemeType = "system"
)

type Theme struct {
	LiteralBase[string]
}

func NewTheme(s string) (*Theme, error) {
	switch s {
	case string(ThemeLight), string(ThemeDark), string(ThemeSystem):
		return &Theme{LiteralBase: LiteralBase[string]{v: s}}, nil
	default:
		return nil, errors.New("invalid theme")
	}
}