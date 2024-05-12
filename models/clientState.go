package models

import (
	"time"

	"github.com/sknutsen/planner/database"
	"github.com/sknutsen/planner/lib"
)

type ClientState struct {
	Plans          []database.Plan
	SelectedPlanId int
	UserProfile    UserProfile
	Week           Week
}

func GetClientState() (ClientState, error) {
	state := ClientState{
		Plans: []database.Plan{
			{
				ID:   1,
				Name: "Plan",
			},
		},
		SelectedPlanId: 1,
		Week: Week{
			ISOWeek: lib.ISOWeek(time.Now()),
			Monday: Day{
				Date:  time.Now(),
				Tasks: []Task{},
			},
			Tuesday: Day{
				Date:  time.Now(),
				Tasks: []Task{},
			},
			Wednesday: Day{
				Date:  time.Now(),
				Tasks: []Task{},
			},
			Thursday: Day{
				Date:  time.Now(),
				Tasks: []Task{},
			},
			Friday: Day{
				Date:  time.Now(),
				Tasks: []Task{},
			},
			Saturday: Day{
				Date:  time.Now(),
				Tasks: []Task{},
			},
			Sunday: Day{
				Date:  time.Now(),
				Tasks: []Task{},
			},
		},
	}

	return state, nil
}
