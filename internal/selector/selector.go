package selector

import (
	"fmt"
	"strings"
)

// Type represents the different types of selectors available.
type Type uint8

const (
	UnknownType Type = iota
	TypeFzy
	TypeFzf
	TypeSkim
)

func TypeFromStr(s string) (Type, error) {
	switch strings.ToLower(s) {
	case "fzy":
		return TypeFzy, nil
	case "fzf":
		return TypeFzf, nil
	case "sk":
		return TypeSkim, nil
	default:
		return UnknownType, fmt.Errorf("Invalid selector '%s'", s)
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
		return nil, fmt.Errorf("Failed to start selector")
	}
}
