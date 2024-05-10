package models

import "time"

type Day struct {
	Date  time.Time
	Tasks []Task
}
