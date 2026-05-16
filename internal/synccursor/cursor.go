package synccursor

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
)

// Keyset is the value encoded in next_cursor (updated_at + row id).
type Keyset struct {
	UpdatedAt string `json:"u"`
	ID        int64  `json:"i"`
}

// Encode returns an opaque URL-safe cursor or empty string for none.
func Encode(k Keyset) string {
	if k.UpdatedAt == "" && k.ID == 0 {
		return ""
	}
	b, err := json.Marshal(k)
	if err != nil {
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

// Decode parses next_cursor; empty s means first page.
func Decode(s string) (Keyset, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Keyset{}, nil
	}
	raw, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return Keyset{}, err
	}
	var k Keyset
	if err := json.Unmarshal(raw, &k); err != nil {
		return Keyset{}, err
	}
	if k.UpdatedAt == "" {
		return Keyset{}, errors.New("invalid cursor")
	}
	return k, nil
}
