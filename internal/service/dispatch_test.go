package service

import (
	"errors"
	"testing"
	"time"

	"fieldops-roster/internal/model"
)

func TestAssignOpenOrdersUpdatesWorkspaceState(t *testing.T) {
	workspace := model.Workspace{
		BusinessDate: mustTime(t, "2026-08-16T09:00:00Z"),
		Technicians: []model.Technician{
			{ID: "t-1", Name: "Aki", Region: "north", Skills: []string{"fiber"}, DailyCapacityMin: 120, Active: true},
		},
		WorkOrders: []model.WorkOrder{
			{ID: "wo-1", Customer: "Clinic A", Region: "north", RequiredSkill: "fiber", DurationMinutes: 60, Priority: model.PriorityUrgent, Status: model.StatusOpen, DueAt: mustTime(t, "2026-08-16T10:00:00Z")},
		},
	}
	dispatch := NewDispatchService(workspace)

	result := dispatch.AssignOpenOrders()

	if len(result.Assignments) != 1 {
		t.Fatalf("assignments = %d, want 1", len(result.Assignments))
	}
	order := dispatch.Workspace.WorkOrders[0]
	if order.Status != model.StatusScheduled || order.AssignedTechID != "t-1" {
		t.Fatalf("order not scheduled correctly: %+v", order)
	}
}

func TestCancelOrderRemovesAssignment(t *testing.T) {
	workspace := model.Workspace{
		Technicians: []model.Technician{
			{ID: "t-1", Name: "Aki", Region: "north", Skills: []string{"fiber"}, DailyCapacityMin: 120, Active: true},
		},
		WorkOrders: []model.WorkOrder{
			{ID: "wo-1", Customer: "Clinic A", Region: "north", RequiredSkill: "fiber", DurationMinutes: 60, Priority: model.PriorityUrgent, Status: model.StatusScheduled, AssignedTechID: "t-1", DueAt: mustTime(t, "2026-08-16T10:00:00Z")},
		},
		Assignments: []model.Assignment{
			{OrderID: "wo-1", TechID: "t-1", ScheduledFor: mustTime(t, "2026-08-16T10:00:00Z"), Minutes: 60},
		},
	}
	dispatch := NewDispatchService(workspace)

	if err := dispatch.CancelOrder("wo-1"); err != nil {
		t.Fatalf("cancel order: %v", err)
	}

	order := dispatch.Workspace.WorkOrders[0]
	if order.Status != model.StatusCancelled || order.AssignedTechID != "" {
		t.Fatalf("order not cancelled correctly: %+v", order)
	}
	if len(dispatch.Workspace.Assignments) != 0 {
		t.Fatalf("assignments = %+v, want empty", dispatch.Workspace.Assignments)
	}
}

func TestCompleteOrderRejectsOpenOrders(t *testing.T) {
	workspace := model.Workspace{
		WorkOrders: []model.WorkOrder{
			{ID: "wo-1", Customer: "Clinic A", Region: "north", RequiredSkill: "fiber", DurationMinutes: 60, Priority: model.PriorityUrgent, Status: model.StatusOpen, DueAt: mustTime(t, "2026-08-16T10:00:00Z")},
		},
	}
	dispatch := NewDispatchService(workspace)

	err := dispatch.CompleteOrder("wo-1", mustTime(t, "2026-08-16T11:00:00Z"))

	if !errors.Is(err, model.ErrInvalidTransition) {
		t.Fatalf("error = %v, want ErrInvalidTransition", err)
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
