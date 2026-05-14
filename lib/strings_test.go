package lib

import "testing"

func TestIsNullOrEmpty_String(t *testing.T) {
	cases := []struct {
		in   interface{}
		want bool
	}{
		{"", true},
		{"   ", true},
		{"\t\n", true},
		{"x", false},
		{" hello ", false},
	}
	for _, tc := range cases {
		if got := IsNullOrEmpty(tc.in); got != tc.want {
			t.Fatalf("IsNullOrEmpty(%#v) = %v want %v", tc.in, got, tc.want)
		}
	}
}

func TestIsNullOrEmpty_Nil(t *testing.T) {
	if !IsNullOrEmpty(nil) {
		t.Fatal("nil should be empty")
	}
}

func TestAsString(t *testing.T) {
	cases := []struct {
		in   interface{}
		want string
	}{
		{nil, ""},
		{"hello", "hello"},
		{42, "42"},
	}
	for _, tc := range cases {
		if got := AsString(tc.in); got != tc.want {
			t.Fatalf("AsString(%#v) = %q want %q", tc.in, got, tc.want)
		}
	}
}
