package store

import (
	"errors"
	"testing"
	"time"

	"fieldops-roster/internal/model"
)

func TestValidateWorkspaceRejectsAssignmentForUnknownOrder(t *testing.T) {
	workspace := model.Workspace{
		Technicians: []model.Technician{
			{ID: "t-1", Name: "Aki", Region: "north", Skills: []string{"fiber"}, DailyCapacityMin: 120, Active: true},
		},
		WorkOrders: []model.WorkOrder{
			{ID: "wo-1", Customer: "Clinic A", Region: "north", RequiredSkill: "fiber", DurationMinutes: 60, Priority: model.PriorityUrgent, Status: model.StatusScheduled, AssignedTechID: "t-1", DueAt: mustTime(t, "2026-08-16T10:00:00Z")},
		},
		Assignments: []model.Assignment{
			{OrderID: "wo-missing", TechID: "t-1", ScheduledFor: mustTime(t, "2026-08-16T10:00:00Z"), Minutes: 60},
		},
	}

	err := ValidateWorkspace(workspace)

	if !errors.Is(err, model.ErrUnknownOrder) {
		t.Fatalf("error = %v, want ErrUnknownOrder", err)
	}
}

func mustTime(t *testing.T, value string) time.Time {
	t.Helper()
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Fatal(err)
	}
	return parsed
}
