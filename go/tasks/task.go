package tasks

import (
	"errors"
	"strings"
	"time"
)

type (
	Entry struct {
		ID             int       `json:"id"`
		AssignedUserID int       `json:"assignedUserId"`
		Title          string    `json:"title"`
		Description    string    `json:"description"`
		Status         Status    `json:"status"`
		CreatedAt      time.Time `json:"createdAt"`
		DueDate        time.Time `json:"dueDate"`
	}

	Status string
)

const (
	StatusPending    Status = "PENDING"
	StatusInProgress Status = "IN_PROGRESS"
	StatusCompleted  Status = "COMPLETED"
)

var Statuses = []Status{StatusPending, StatusInProgress, StatusCompleted}

func (e Entry) Pending() bool {
	return e.Status == StatusPending
}

func (e Entry) InProgress() bool {
	return e.Status == StatusInProgress
}

func (e Entry) Completed() bool {
	return e.Status == StatusCompleted
}

func (e Entry) Validate() error {
	if e.Title == "" {
		return errors.New("missing title")
	}
	if e.AssignedUserID <= 0 {
		return errors.New("invalid assigned user")
	}
	if time.Now().After(e.DueDate) {
		return errors.New("due date expired")
	}
	return nil
}

func ParseStatus(str string) (Status, error) {
	for _, s := range Statuses {
		if strings.EqualFold(str, string(s)) {
			return s, nil
		}
	}
	return "", errors.New("invalid status")
}
