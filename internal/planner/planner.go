package planner

import (
	"sort"
	"time"

	"fieldops-roster/internal/model"
)

type Candidate struct {
	Technician model.Technician
	UsedMin    int
}

type PlanResult struct {
	Assignments []model.Assignment
	Unassigned  []model.WorkOrder
}

func Plan(workspace model.Workspace) PlanResult {
	used := capacityByTechnician(workspace.Assignments)
	orders := pendingOrders(workspace.WorkOrders)
	sort.SliceStable(orders, func(i, j int) bool {
		if orders[i].Priority.Rank() != orders[j].Priority.Rank() {
			return orders[i].Priority.Rank() > orders[j].Priority.Rank()
		}
		if !orders[i].DueAt.Equal(orders[j].DueAt) {
			return orders[i].DueAt.Before(orders[j].DueAt)
		}
		return orders[i].ID < orders[j].ID
	})

	var result PlanResult
	for _, order := range orders {
		tech, ok := bestTechnician(workspace.Technicians, used, order)
		if !ok {
			result.Unassigned = append(result.Unassigned, order)
			continue
		}
		used[tech.ID] += order.DurationMinutes
		result.Assignments = append(result.Assignments, model.Assignment{
			OrderID:      order.ID,
			TechID:       tech.ID,
			ScheduledFor: scheduleDay(workspace.BusinessDate, order.DueAt),
			Minutes:      order.DurationMinutes,
		})
	}
	return result
}

func capacityByTechnician(assignments []model.Assignment) map[string]int {
	used := make(map[string]int)
	for _, assignment := range assignments {
		used[assignment.TechID] += assignment.Minutes
	}
	return used
}

func pendingOrders(orders []model.WorkOrder) []model.WorkOrder {
	pending := make([]model.WorkOrder, 0, len(orders))
	for _, order := range orders {
		if order.Status == model.StatusOpen && order.AssignedTechID == "" {
			pending = append(pending, order)
		}
	}
	return pending
}

func bestTechnician(techs []model.Technician, used map[string]int, order model.WorkOrder) (model.Technician, bool) {
	candidates := make([]Candidate, 0, len(techs))
	for _, tech := range techs {
		if !tech.Active || tech.Region != order.Region || !tech.HasSkill(order.RequiredSkill) {
			continue
		}
		if used[tech.ID]+order.DurationMinutes > tech.DailyCapacityMin {
			continue
		}
		candidates = append(candidates, Candidate{Technician: tech, UsedMin: used[tech.ID]})
	}
	if len(candidates) == 0 {
		return model.Technician{}, false
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].UsedMin != candidates[j].UsedMin {
			return candidates[i].UsedMin < candidates[j].UsedMin
		}
		return candidates[i].Technician.ID < candidates[j].Technician.ID
	})
	return candidates[0].Technician, true
}

func scheduleDay(businessDate, due time.Time) time.Time {
	if businessDate.IsZero() {
		return due.Truncate(24 * time.Hour)
	}
	if due.Before(businessDate) {
		return businessDate
	}
	return due
}
