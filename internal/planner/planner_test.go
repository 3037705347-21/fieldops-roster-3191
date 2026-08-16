package planner

import (
	"testing"
	"time"

	"fieldops-roster/internal/model"
)

func TestPlanPrioritizesUrgentDueOrdersAndBalancesCapacity(t *testing.T) {
	workspace := model.Workspace{
		BusinessDate: mustTime(t, "2026-08-16T09:00:00Z"),
		Technicians: []model.Technician{
			{ID: "t-1", Name: "Aki", Region: "north", Skills: []string{"fiber"}, DailyCapacityMin: 180, Active: true},
			{ID: "t-2", Name: "Mina", Region: "north", Skills: []string{"fiber"}, DailyCapacityMin: 180, Active: true},
		},
		WorkOrders: []model.WorkOrder{
			{ID: "wo-low", Customer: "Cafe Blue", Region: "north", RequiredSkill: "fiber", DurationMinutes: 90, Priority: model.PriorityLow, Status: model.StatusOpen, DueAt: mustTime(t, "2026-08-18T10:00:00Z")},
			{ID: "wo-urgent", Customer: "Clinic A", Region: "north", RequiredSkill: "fiber", DurationMinutes: 90, Priority: model.PriorityUrgent, Status: model.StatusOpen, DueAt: mustTime(t, "2026-08-17T10:00:00Z")},
			{ID: "wo-standard", Customer: "Office 9", Region: "north", RequiredSkill: "fiber", DurationMinutes: 90, Priority: model.PriorityStandard, Status: model.StatusOpen, DueAt: mustTime(t, "2026-08-16T16:00:00Z")},
		},
	}

	result := Plan(workspace)

	if len(result.Assignments) != 3 {
		t.Fatalf("assignments = %d, want 3", len(result.Assignments))
	}
	if result.Assignments[0].OrderID != "wo-urgent" {
		t.Fatalf("first assignment = %s, want urgent order first", result.Assignments[0].OrderID)
	}
	if result.Assignments[0].TechID != "t-1" || result.Assignments[1].TechID != "t-2" {
		t.Fatalf("expected planner to balance first two orders across technicians: %+v", result.Assignments[:2])
	}
}

func TestPlanLeavesOrdersUnassignedWhenSkillOrCapacityIsMissing(t *testing.T) {
	workspace := model.Workspace{
		BusinessDate: mustTime(t, "2026-08-16T09:00:00Z"),
		Technicians: []model.Technician{
			{ID: "t-1", Name: "Aki", Region: "north", Skills: []string{"fiber"}, DailyCapacityMin: 100, Active: true},
		},
		WorkOrders: []model.WorkOrder{
			{ID: "wo-fit", Customer: "Clinic A", Region: "north", RequiredSkill: "fiber", DurationMinutes: 80, Priority: model.PriorityUrgent, Status: model.StatusOpen, DueAt: mustTime(t, "2026-08-16T10:00:00Z")},
			{ID: "wo-too-large", Customer: "Office 9", Region: "north", RequiredSkill: "fiber", DurationMinutes: 40, Priority: model.PriorityStandard, Status: model.StatusOpen, DueAt: mustTime(t, "2026-08-16T11:00:00Z")},
			{ID: "wo-skill", Customer: "Cafe Blue", Region: "north", RequiredSkill: "hvac", DurationMinutes: 30, Priority: model.PriorityLow, Status: model.StatusOpen, DueAt: mustTime(t, "2026-08-16T12:00:00Z")},
		},
	}

	result := Plan(workspace)

	if len(result.Assignments) != 1 {
		t.Fatalf("assignments = %d, want 1", len(result.Assignments))
	}
	if len(result.Unassigned) != 2 {
		t.Fatalf("unassigned = %d, want 2", len(result.Unassigned))
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
