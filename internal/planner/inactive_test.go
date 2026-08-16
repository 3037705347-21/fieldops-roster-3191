package planner

import (
	"testing"
	"time"

	"fieldops-roster/internal/model"
)

func TestPlanSkipsInactiveTechniciansEvenWhenTheyMatchTheSkill(t *testing.T) {
	workspace := model.Workspace{
		BusinessDate: mustInactiveTime(t, "2026-08-16T09:00:00Z"),
		Technicians: []model.Technician{
			{ID: "t-inactive", Name: "Inactive Lead", Region: "north", Skills: []string{"fiber"}, DailyCapacityMin: 300, Active: false},
			{ID: "t-active", Name: "Active Backup", Region: "north", Skills: []string{"router"}, DailyCapacityMin: 300, Active: true},
		},
		WorkOrders: []model.WorkOrder{
			{ID: "wo-fiber", Customer: "Clinic A", Region: "north", RequiredSkill: "fiber", DurationMinutes: 90, Priority: model.PriorityUrgent, Status: model.StatusOpen, DueAt: mustInactiveTime(t, "2026-08-16T12:00:00Z")},
		},
	}

	result := Plan(workspace)

	if len(result.Assignments) != 0 {
		t.Fatalf("inactive technician received work: %+v", result.Assignments)
	}
	if len(result.Unassigned) != 1 {
		t.Fatalf("unassigned = %d, want 1", len(result.Unassigned))
	}
}

func mustInactiveTime(t *testing.T, value string) time.Time {
	t.Helper()
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Fatal(err)
	}
	return parsed
}
