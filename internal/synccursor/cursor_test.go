package synccursor_test

import (
	"testing"

	"github.com/sknutsen/planner/internal/synccursor"
)

func TestEncodeDecodeRoundTrip(t *testing.T) {
	k := synccursor.Keyset{UpdatedAt: "2026-05-15T10:00:00.000000Z", ID: 42}
	s := synccursor.Encode(k)
	got, err := synccursor.Decode(s)
	if err != nil {
		t.Fatal(err)
	}
	if got != k {
		t.Fatalf("got %+v want %+v", got, k)
	}
}

func TestDecodeEmpty(t *testing.T) {
	k, err := synccursor.Decode("")
	if err != nil || (k.ID != 0 || k.UpdatedAt != "") {
		t.Fatalf("empty decode: %+v err=%v", k, err)
	}
}
