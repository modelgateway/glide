package routing

import (
	"errors"

	"glide/pkg/providers"
)

var ErrNoHealthyModels = errors.New("no healthy models found")

// Strategy defines supported routing strategies for language routers
type Strategy string

type LangModelRouting interface {
	Iterator() LangModelIterator
}

type LangModelIterator interface {
	Next() (*providers.LangModel, error)
}
