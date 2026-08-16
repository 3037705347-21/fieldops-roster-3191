package model

import (
	"errors"
	"fmt"
	"time"
)

type Priority string

const (
	PriorityLow      Priority = "low"
	PriorityStandard Priority = "standard"
	PriorityUrgent   Priority = "urgent"
)

type WorkOrderStatus string

const (
	StatusOpen       WorkOrderStatus = "open"
	StatusScheduled  WorkOrderStatus = "scheduled"
	StatusInProgress WorkOrderStatus = "in_progress"
	StatusCompleted  WorkOrderStatus = "completed"
	StatusCancelled  WorkOrderStatus = "cancelled"
)

type Technician struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Region           string   `json:"region"`
	Skills           []string `json:"skills"`
	DailyCapacityMin int      `json:"daily_capacity_min"`
	Active           bool     `json:"active"`
}

type WorkOrder struct {
	ID              string          `json:"id"`
	Customer        string          `json:"customer"`
	Region          string          `json:"region"`
	RequiredSkill   string          `json:"required_skill"`
	DurationMinutes int             `json:"duration_minutes"`
	Priority        Priority        `json:"priority"`
	Status          WorkOrderStatus `json:"status"`
	DueAt           time.Time       `json:"due_at"`
	AssignedTechID  string          `json:"assigned_tech_id,omitempty"`
	CompletedAt     *time.Time      `json:"completed_at,omitempty"`
}

type Assignment struct {
	OrderID      string    `json:"order_id"`
	TechID       string    `json:"tech_id"`
	ScheduledFor time.Time `json:"scheduled_for"`
	Minutes      int       `json:"minutes"`
}

type Workspace struct {
	GeneratedAt  time.Time    `json:"generated_at"`
	Technicians  []Technician `json:"technicians"`
	WorkOrders   []WorkOrder  `json:"work_orders"`
	Assignments  []Assignment `json:"assignments"`
	BusinessDate time.Time    `json:"business_date"`
}

var (
	ErrUnknownTechnician = errors.New("unknown technician")
	ErrUnknownOrder      = errors.New("unknown work order")
	ErrInvalidTransition = errors.New("invalid work order transition")
	ErrCapacityExceeded  = errors.New("technician capacity exceeded")
)

func (p Priority) Rank() int {
	switch p {
	case PriorityUrgent:
		return 0
	case PriorityStandard:
		return 1
	case PriorityLow:
		return 2
	default:
		return 3
	}
}

func (o WorkOrder) Validate() error {
	if o.ID == "" {
		return fmt.Errorf("work order id is required")
	}
	if o.Customer == "" {
		return fmt.Errorf("work order %s customer is required", o.ID)
	}
	if o.Region == "" {
		return fmt.Errorf("work order %s region is required", o.ID)
	}
	if o.RequiredSkill == "" {
		return fmt.Errorf("work order %s required skill is required", o.ID)
	}
	if o.DurationMinutes <= 0 {
		return fmt.Errorf("work order %s duration must be positive", o.ID)
	}
	if o.Status == "" {
		return fmt.Errorf("work order %s status is required", o.ID)
	}
	return nil
}

func (t Technician) HasSkill(skill string) bool {
	for _, candidate := range t.Skills {
		if candidate == skill {
			return true
		}
	}
	return false
}

func (t Technician) Validate() error {
	if t.ID == "" {
		return fmt.Errorf("technician id is required")
	}
	if t.Name == "" {
		return fmt.Errorf("technician %s name is required", t.ID)
	}
	if t.Region == "" {
		return fmt.Errorf("technician %s region is required", t.ID)
	}
	if t.DailyCapacityMin <= 0 {
		return fmt.Errorf("technician %s capacity must be positive", t.ID)
	}
	if len(t.Skills) == 0 {
		return fmt.Errorf("technician %s needs at least one skill", t.ID)
	}
	return nil
}

func CanTransition(from, to WorkOrderStatus) bool {
	switch from {
	case StatusOpen:
		return to == StatusScheduled || to == StatusCompleted || to == StatusCancelled
	case StatusScheduled:
		return to == StatusInProgress || to == StatusCancelled || to == StatusOpen
	case StatusInProgress:
		return to == StatusCompleted || to == StatusScheduled
	case StatusCompleted, StatusCancelled:
		return false
	default:
		return false
	}
}
