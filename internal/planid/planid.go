// Package planid resolves the active plan ID for list UIs from the user's plans
// and an optional requested plan ID from the route.
package planid

import "github.com/sknutsen/planner/database"

// Selected returns the plan ID to use: requestedID when it matches a plan, the
// first plan's ID when requestedID is 0, or 0 when plans is empty or there is no match.
func Selected(plans []database.Plan, requestedID int) int {
	for _, p := range plans {
		if requestedID == 0 || requestedID == int(p.ID) {
			return int(p.ID)
		}
	}
	return 0
}
