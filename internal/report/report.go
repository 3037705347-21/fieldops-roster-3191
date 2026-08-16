package report

import (
	"sort"
	"time"

	"fieldops-roster/internal/model"
)

type RegionSummary struct {
	Region             string
	Open               int
	Scheduled          int
	Completed          int
	Cancelled          int
	OverdueOpen        int
	AssignedMinutes    int
	AvailableTechCount int
}

type Summary struct {
	GeneratedAt time.Time
	Regions     []RegionSummary
}

func BuildSummary(workspace model.Workspace, now time.Time) Summary {
	byRegion := make(map[string]*RegionSummary)
	for _, tech := range workspace.Technicians {
		summary := ensureRegion(byRegion, tech.Region)
		if tech.Active {
			summary.AvailableTechCount++
		}
	}
	for _, order := range workspace.WorkOrders {
		summary := ensureRegion(byRegion, order.Region)
		switch order.Status {
		case model.StatusOpen:
			summary.Open++
			if order.DueAt.Before(now) {
				summary.OverdueOpen++
			}
		case model.StatusScheduled:
			summary.Scheduled++
		case model.StatusCompleted:
			summary.Completed++
		case model.StatusCancelled:
			summary.Cancelled++
		}
	}
	for _, assignment := range workspace.Assignments {
		for _, order := range workspace.WorkOrders {
			if order.ID == assignment.OrderID {
				ensureRegion(byRegion, order.Region).AssignedMinutes += assignment.Minutes
				break
			}
		}
	}

	regions := make([]RegionSummary, 0, len(byRegion))
	for _, summary := range byRegion {
		regions = append(regions, *summary)
	}
	sort.SliceStable(regions, func(i, j int) bool {
		return regions[i].Region < regions[j].Region
	})
	return Summary{GeneratedAt: now, Regions: regions}
}

func ensureRegion(regions map[string]*RegionSummary, region string) *RegionSummary {
	if region == "" {
		region = "unassigned"
	}
	if _, exists := regions[region]; !exists {
		regions[region] = &RegionSummary{Region: region}
	}
	return regions[region]
}
