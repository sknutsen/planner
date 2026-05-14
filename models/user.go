package models

import (
	"encoding/json"
	"fmt"
)

type UserProfile struct {
	UserId    string
	Name      string
	Nickname  string
	Picture   string
	Iat       float64
	Exp       float64
	UpdatedAt string
	Admin     bool
}

func stringClaim(profile map[string]interface{}, key string) string {
	v, ok := profile[key]
	if !ok || v == nil {
		return ""
	}
	switch s := v.(type) {
	case string:
		return s
	default:
		return fmt.Sprint(s)
	}
}

func floatClaim(profile map[string]interface{}, key string) float64 {
	v, ok := profile[key]
	if !ok || v == nil {
		return 0
	}
	switch n := v.(type) {
	case float64:
		return n
	case float32:
		return float64(n)
	case int:
		return float64(n)
	case int64:
		return float64(n)
	case json.Number:
		f, err := n.Float64()
		if err != nil {
			return 0
		}
		return f
	default:
		return 0
	}
}

// GetUserProfile maps an OIDC / Auth0 profile map into UserProfile. Missing or
// mistyped claims become zero values instead of panicking.
func GetUserProfile(profile map[string]interface{}) UserProfile {
	if profile == nil {
		return UserProfile{}
	}
	return UserProfile{
		UserId:    stringClaim(profile, "sub"),
		Name:      stringClaim(profile, "name"),
		Picture:   stringClaim(profile, "picture"),
		Nickname:  stringClaim(profile, "nickname"),
		Iat:       floatClaim(profile, "iat"),
		Exp:       floatClaim(profile, "exp"),
		UpdatedAt: stringClaim(profile, "updated_at"),
	}
}
