package service

import (
	"fmt"
	"time"

	"fieldops-roster/internal/model"
	"fieldops-roster/internal/planner"
)

type DispatchService struct {
	Workspace model.Workspace
}

func NewDispatchService(workspace model.Workspace) DispatchService {
	return DispatchService{Workspace: workspace}
}

func (s *DispatchService) AssignOpenOrders() planner.PlanResult {
	result := planner.Plan(s.Workspace)
	for _, assignment := range result.Assignments {
		s.Workspace.Assignments = append(s.Workspace.Assignments, assignment)
		for i := range s.Workspace.WorkOrders {
			if s.Workspace.WorkOrders[i].ID == assignment.OrderID {
				s.Workspace.WorkOrders[i].Status = model.StatusScheduled
				s.Workspace.WorkOrders[i].AssignedTechID = assignment.TechID
				break
			}
		}
	}
	return result
}

func (s *DispatchService) CompleteOrder(orderID string, completedAt time.Time) error {
	order, idx, ok := s.findOrder(orderID)
	if !ok {
		return fmt.Errorf("%w: %s", model.ErrUnknownOrder, orderID)
	}
	if !model.CanTransition(order.Status, model.StatusCompleted) {
		return fmt.Errorf("%w: %s to %s", model.ErrInvalidTransition, order.Status, model.StatusCompleted)
	}
	s.Workspace.WorkOrders[idx].Status = model.StatusCompleted
	s.Workspace.WorkOrders[idx].CompletedAt = &completedAt
	return nil
}

func (s *DispatchService) RequeueOrder(orderID string) error {
	order, idx, ok := s.findOrder(orderID)
	if !ok {
		return fmt.Errorf("%w: %s", model.ErrUnknownOrder, orderID)
	}
	if !model.CanTransition(order.Status, model.StatusOpen) {
		return fmt.Errorf("%w: %s to %s", model.ErrInvalidTransition, order.Status, model.StatusOpen)
	}
	s.Workspace.WorkOrders[idx].Status = model.StatusOpen
	s.Workspace.WorkOrders[idx].AssignedTechID = ""
	s.removeAssignment(orderID)
	return nil
}

func (s *DispatchService) CancelOrder(orderID string) error {
	order, idx, ok := s.findOrder(orderID)
	if !ok {
		return fmt.Errorf("%w: %s", model.ErrUnknownOrder, orderID)
	}
	if !model.CanTransition(order.Status, model.StatusCancelled) {
		return fmt.Errorf("%w: %s to %s", model.ErrInvalidTransition, order.Status, model.StatusCancelled)
	}
	s.Workspace.WorkOrders[idx].Status = model.StatusCancelled
	s.Workspace.WorkOrders[idx].AssignedTechID = ""
	return nil
}

func (s *DispatchService) findOrder(orderID string) (model.WorkOrder, int, bool) {
	for idx, order := range s.Workspace.WorkOrders {
		if order.ID == orderID {
			return order, idx, true
		}
	}
	return model.WorkOrder{}, -1, false
}

func (s *DispatchService) removeAssignment(orderID string) {
	filtered := s.Workspace.Assignments[:0]
	for _, assignment := range s.Workspace.Assignments {
		if assignment.OrderID != orderID {
			filtered = append(filtered, assignment)
		}
	}
	s.Workspace.Assignments = filtered
}
