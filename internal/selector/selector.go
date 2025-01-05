package selector

import "strings"

// Type represents the different types of selectors available.
type Type uint8

const (
	UnknownType Type = iota
	TypeFzy
	TypeFzf
	TypeSkim
)

func TypeFromStr(s string) Type {
	switch strings.ToLower(s) {
	case "fzy":
		return TypeFzy
	case "fzf":
		return TypeFzf
	case "sk":
		return TypeSkim
	default:
		return UnknownType
	}
}

// Displays a series of options for user selection.
type Selector interface {
	Run(inputChan chan string) (string, error)
}

// New creates a new Selector instance based on the provided selector type and options.
func New(t Type) (Selector, error) {
	switch t {
	case TypeFzf:
		return NewCmd("fzf"), nil
	case TypeFzy:
		return NewCmd("fzy"), nil
	case TypeSkim:
		return NewCmd("sk"), nil
	default:
		return NewCmd("fzf"), nil
	}
}
