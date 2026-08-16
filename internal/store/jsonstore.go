package store

import (
	"encoding/json"
	"fmt"
	"os"

	"fieldops-roster/internal/model"
)

type JSONStore struct {
	Path string
}

func NewJSONStore(path string) JSONStore {
	return JSONStore{Path: path}
}

func (s JSONStore) Load() (model.Workspace, error) {
	data, err := os.ReadFile(s.Path)
	if err != nil {
		return model.Workspace{}, fmt.Errorf("read workspace: %w", err)
	}
	var workspace model.Workspace
	if err := json.Unmarshal(data, &workspace); err != nil {
		return model.Workspace{}, fmt.Errorf("decode workspace: %w", err)
	}
	if err := ValidateWorkspace(workspace); err != nil {
		return model.Workspace{}, err
	}
	return workspace, nil
}

func (s JSONStore) Save(workspace model.Workspace) error {
	if err := ValidateWorkspace(workspace); err != nil {
		return err
	}
	data, err := json.MarshalIndent(workspace, "", "  ")
	if err != nil {
		return fmt.Errorf("encode workspace: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(s.Path, data, 0o644); err != nil {
		return fmt.Errorf("write workspace: %w", err)
	}
	return nil
}

func ValidateWorkspace(workspace model.Workspace) error {
	techIDs := make(map[string]struct{}, len(workspace.Technicians))
	for _, tech := range workspace.Technicians {
		if err := tech.Validate(); err != nil {
			return err
		}
		if _, exists := techIDs[tech.ID]; exists {
			return fmt.Errorf("duplicate technician %s", tech.ID)
		}
		techIDs[tech.ID] = struct{}{}
	}

	orderIDs := make(map[string]struct{}, len(workspace.WorkOrders))
	for _, order := range workspace.WorkOrders {
		if err := order.Validate(); err != nil {
			return err
		}
		if _, exists := orderIDs[order.ID]; exists {
			return fmt.Errorf("duplicate work order %s", order.ID)
		}
		if order.AssignedTechID != "" {
			if _, exists := techIDs[order.AssignedTechID]; !exists {
				return fmt.Errorf("%w: %s", model.ErrUnknownTechnician, order.AssignedTechID)
			}
		}
		orderIDs[order.ID] = struct{}{}
	}

	for _, assignment := range workspace.Assignments {
		if _, exists := techIDs[assignment.TechID]; !exists {
			return fmt.Errorf("%w: %s", model.ErrUnknownTechnician, assignment.TechID)
		}
		if _, exists := orderIDs[assignment.OrderID]; !exists {
			return fmt.Errorf("%w: %s", model.ErrUnknownOrder, assignment.OrderID)
		}
		if assignment.Minutes <= 0 {
			return fmt.Errorf("assignment %s minutes must be positive", assignment.OrderID)
		}
	}
	return nil
}
