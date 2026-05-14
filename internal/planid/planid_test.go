package planid

import (
	"testing"

	"github.com/sknutsen/planner/database"
)

func TestSelected(t *testing.T) {
	plans := []database.Plan{
		{ID: 10, Name: "A"},
		{ID: 20, Name: "B"},
	}
	cases := []struct {
		name string
		req  int
		want int
	}{
		{"zero picks first", 0, 10},
		{"match first", 10, 10},
		{"match second", 20, 20},
		{"no match", 99, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Selected(plans, tc.req); got != tc.want {
				t.Fatalf("Selected(_, %d) = %d want %d", tc.req, got, tc.want)
			}
		})
	}
}

func TestSelected_EmptyPlans(t *testing.T) {
	if got := Selected(nil, 0); got != 0 {
		t.Fatalf("got %d", got)
	}
	if got := Selected([]database.Plan{}, 10); got != 0 {
		t.Fatalf("got %d", got)
	}
}
