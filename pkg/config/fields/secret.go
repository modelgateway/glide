package fields

import (
	"encoding"
)

// Secret is a string that is marshaled in an opaque way, so we are not leaking sensitive information
type Secret string

const maskedSecret = "[REDACTED]"

var _ encoding.TextMarshaler = Secret("")

// MarshalText marshals the secret as `[REDACTED]`.
func (s Secret) MarshalText() ([]byte, error) {
	return []byte(maskedSecret), nil
}
