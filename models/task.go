package models

import "time"

type Task struct {
	Id          int
	Date        time.Time
	Title       string
	Subtitle    string
	Description string
}
