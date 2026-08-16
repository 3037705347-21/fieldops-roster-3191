package report

import (
	"testing"
	"time"

	"fieldops-roster/internal/model"
)

func TestBuildSummaryGroupsRegionsAndCountsOverdueOpenWork(t *testing.T) {
	now := mustTime(t, "2026-08-16T12:00:00Z")
	workspace := model.Workspace{
		Technicians: []model.Technician{
			{ID: "t-1", Name: "Aki", Region: "north", Skills: []string{"fiber"}, DailyCapacityMin: 120, Active: true},
			{ID: "t-2", Name: "Mina", Region: "south", Skills: []string{"hvac"}, DailyCapacityMin: 120, Active: false},
		},
		WorkOrders: []model.WorkOrder{
			{ID: "wo-1", Customer: "Clinic A", Region: "north", RequiredSkill: "fiber", DurationMinutes: 60, Priority: model.PriorityUrgent, Status: model.StatusOpen, DueAt: mustTime(t, "2026-08-16T10:00:00Z")},
			{ID: "wo-2", Customer: "Office 9", Region: "north", RequiredSkill: "fiber", DurationMinutes: 90, Priority: model.PriorityStandard, Status: model.StatusScheduled, DueAt: mustTime(t, "2026-08-16T16:00:00Z")},
			{ID: "wo-3", Customer: "Cafe Blue", Region: "south", RequiredSkill: "hvac", DurationMinutes: 30, Priority: model.PriorityLow, Status: model.StatusCompleted, DueAt: mustTime(t, "2026-08-15T16:00:00Z")},
		},
		Assignments: []model.Assignment{
			{OrderID: "wo-2", TechID: "t-1", ScheduledFor: now, Minutes: 90},
		},
	}

	summary := BuildSummary(workspace, now)

	if len(summary.Regions) != 2 {
		t.Fatalf("regions = %d, want 2", len(summary.Regions))
	}
	north := summary.Regions[0]
	if north.Region != "north" || north.Open != 1 || north.Scheduled != 1 || north.OverdueOpen != 1 || north.AssignedMinutes != 90 || north.AvailableTechCount != 1 {
		t.Fatalf("north summary mismatch: %+v", north)
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
