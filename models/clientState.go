package models

import (
	"time"

	"github.com/sknutsen/planner/database"
	"github.com/sknutsen/planner/lib"
)

type ClientState struct {
	BaseRoute      string
	Plans          []database.Plan
	SelectedPlanId int
	UserProfile    UserProfile
}

type HistoryState struct {
	State ClientState
	Tasks []Task
}

type ResourcesState struct {
	State ClientState
	Tasks []Task
}

type TemplatesState struct {
	State     ClientState
	Templates []Template
}

type UserState struct {
	State ClientState
}

type WeekState struct {
	State ClientState
	Week  Week
}

func GetClientState() (ClientState, error) {
	state := ClientState{
		Plans:          []database.Plan{},
		SelectedPlanId: 0,
	}

	return state, nil
}

func GetHistoryState() (HistoryState, error) {
	state := HistoryState{
		State: ClientState{
			Plans:          []database.Plan{},
			SelectedPlanId: 0,
		},
		Tasks: []Task{},
	}

	return state, nil
}

func GetResourcesState() (ResourcesState, error) {
	state := ResourcesState{
		State: ClientState{
			Plans:          []database.Plan{},
			SelectedPlanId: 0,
		},
		Tasks: []Task{},
	}

	return state, nil
}

func GetTemplatesState() (TemplatesState, error) {
	state := TemplatesState{
		State: ClientState{
			Plans:          []database.Plan{},
			SelectedPlanId: 0,
		},
		Templates: []Template{},
	}

	return state, nil
}

func GetUserState() (UserState, error) {
	state := UserState{
		State: ClientState{
			Plans:          []database.Plan{},
			SelectedPlanId: 0,
		},
	}

	return state, nil
}

func GetWeekState() (WeekState, error) {
	state := WeekState{
		State: ClientState{
			Plans:          []database.Plan{},
			SelectedPlanId: 0,
		},
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
