package models

import (
	"encoding/json"
	"math"
	"testing"
)

func TestGetUserProfile_Nil(t *testing.T) {
	got := GetUserProfile(nil)
	if got != (UserProfile{}) {
		t.Fatalf("expected zero UserProfile, got %+v", got)
	}
}

func TestGetUserProfile_FullStringClaims(t *testing.T) {
	in := map[string]interface{}{
		"sub":        "auth0|123",
		"name":       "Ada Lovelace",
		"nickname":   "ada",
		"picture":    "https://example.com/p.png",
		"iat":        float64(1000),
		"exp":        float64(2000),
		"updated_at": "2025-01-01",
	}
	got := GetUserProfile(in)
	if got.UserId != "auth0|123" || got.Name != "Ada Lovelace" || got.Nickname != "ada" || got.Picture != "https://example.com/p.png" {
		t.Fatalf("string fields: got %+v", got)
	}
	if got.Iat != 1000 || got.Exp != 2000 {
		t.Fatalf("float fields: got %+v", got)
	}
	if got.UpdatedAt != "2025-01-01" {
		t.Fatalf("UpdatedAt: got %q", got.UpdatedAt)
	}
}

func TestGetUserProfile_MissingClaimsAreEmpty(t *testing.T) {
	got := GetUserProfile(map[string]interface{}{})
	if got.UserId != "" || got.Name != "" {
		t.Fatalf("expected empty strings, got %+v", got)
	}
	if got.Iat != 0 || got.Exp != 0 {
		t.Fatalf("expected zero floats, got %+v", got)
	}
}

func TestGetUserProfile_NumericIatExp(t *testing.T) {
	cases := []struct {
		name string
		iat  interface{}
		exp  float64
	}{
		{"float64", float64(3.5), 3.5},
		{"float32", float32(2), 2},
		{"int", int(7), 7},
		{"int64", int64(9), 9},
		{"json.Number", json.Number("1.25"), 1.25},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := GetUserProfile(map[string]interface{}{"iat": tc.iat})
			if math.Abs(got.Iat-tc.exp) > 1e-9 {
				t.Fatalf("Iat: want %v got %v", tc.exp, got.Iat)
			}
		})
	}
}

func TestGetUserProfile_NonStringUsesSprint(t *testing.T) {
	got := GetUserProfile(map[string]interface{}{"sub": 42})
	if got.UserId != "42" {
		t.Fatalf("sub: want %q got %q", "42", got.UserId)
	}
}
