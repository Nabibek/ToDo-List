package domain

import "time"

type ToDo struct {
	Id            string    `json:"id"`
	Todo          string    `json:"todo"`
	Message       string    `json:"message"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	Deadline      time.Time `json:"deadline"`
	SystemMessage string    `json:"systemMessage"`
	CompletedAt   time.Time `json:"completedAt"`
	Complete      bool      `json:"complete"`
}
