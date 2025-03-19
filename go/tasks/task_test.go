package tasks

import (
	"testing"
	"time"
)

func TestEntry_Status(t *testing.T) {
	cases := []struct {
		name         string
		entry        Entry
		isCompleted  bool
		isInProgress bool
		isPending    bool
	}{
		{
			name:         "Completed",
			entry:        Entry{Status: StatusCompleted},
			isCompleted:  true,
			isInProgress: false,
			isPending:    false,
		},
		{
			name:         "In Progress",
			entry:        Entry{Status: StatusInProgress},
			isCompleted:  false,
			isInProgress: true,
			isPending:    false,
		},
		{
			name:         "Pending",
			entry:        Entry{Status: StatusPending},
			isCompleted:  false,
			isInProgress: false,
			isPending:    true,
		},
	}

	toBeOrNotToBe := func(v bool) string {
		if v {
			return "should be"
		}
		return "should not be"
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.entry.Completed() != tc.isCompleted {
				t.Errorf("Entry %s completed", toBeOrNotToBe(tc.isCompleted))
			}
			if tc.entry.InProgress() != tc.isInProgress {
				t.Errorf("Entry %s be in progress", toBeOrNotToBe(tc.isInProgress))
			}
			if tc.entry.Pending() != tc.isPending {
				t.Errorf("Entry %s pending", toBeOrNotToBe(tc.isPending))
			}
		})
	}
}

func TestEntry_Validate(t *testing.T) {
	cases := []struct {
		name    string
		entry   Entry
		wantErr bool
	}{
		{
			name:    "Valid",
			entry:   Entry{Title: "Valid Title", AssignedUserID: 1, DueDate: time.Now().Add(time.Hour)},
			wantErr: false,
		}, {
			name:    "Missing Title",
			entry:   Entry{AssignedUserID: 1},
			wantErr: true,
		},
		{
			name:    "Negative Assigned User ID",
			entry:   Entry{Title: "Valid Title", AssignedUserID: -1, DueDate: time.Now().Add(-time.Hour)},
			wantErr: true,
		},
		{
			name:    "Zero Assigned User ID",
			entry:   Entry{Title: "Valid Title", AssignedUserID: 0, DueDate: time.Now().Add(-time.Hour)},
			wantErr: true,
		},
		{
			name:    "Due Date Expired",
			entry:   Entry{Title: "Valid Title", AssignedUserID: 1, DueDate: time.Now().Add(time.Hour)},
			wantErr: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.entry.Validate()
			if tc.wantErr && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	}
}
