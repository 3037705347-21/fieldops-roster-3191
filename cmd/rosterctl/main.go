package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"fieldops-roster/internal/report"
	"fieldops-roster/internal/service"
	"fieldops-roster/internal/store"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: rosterctl <validate|plan|complete|cancel|report> -file workspace.json")
	}
	switch args[0] {
	case "validate":
		fs := flag.NewFlagSet("validate", flag.ContinueOnError)
		path := fs.String("file", "workspace.json", "workspace JSON path")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		_, err := store.NewJSONStore(*path).Load()
		return err
	case "plan":
		fs := flag.NewFlagSet("plan", flag.ContinueOnError)
		path := fs.String("file", "workspace.json", "workspace JSON path")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		st := store.NewJSONStore(*path)
		workspace, err := st.Load()
		if err != nil {
			return err
		}
		dispatch := service.NewDispatchService(workspace)
		result := dispatch.AssignOpenOrders()
		if err := st.Save(dispatch.Workspace); err != nil {
			return err
		}
		fmt.Printf("assigned=%d unassigned=%d\n", len(result.Assignments), len(result.Unassigned))
		return nil
	case "complete":
		fs := flag.NewFlagSet("complete", flag.ContinueOnError)
		path := fs.String("file", "workspace.json", "workspace JSON path")
		orderID := fs.String("order", "", "work order id")
		when := fs.String("at", time.Now().UTC().Format(time.RFC3339), "completion timestamp")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		if *orderID == "" {
			return fmt.Errorf("complete requires -order")
		}
		completedAt, err := time.Parse(time.RFC3339, *when)
		if err != nil {
			return fmt.Errorf("parse -at: %w", err)
		}
		st := store.NewJSONStore(*path)
		workspace, err := st.Load()
		if err != nil {
			return err
		}
		dispatch := service.NewDispatchService(workspace)
		if err := dispatch.CompleteOrder(*orderID, completedAt); err != nil {
			return err
		}
		return st.Save(dispatch.Workspace)
	case "cancel":
		fs := flag.NewFlagSet("cancel", flag.ContinueOnError)
		path := fs.String("file", "workspace.json", "workspace JSON path")
		orderID := fs.String("order", "", "work order id")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		if *orderID == "" {
			return fmt.Errorf("cancel requires -order")
		}
		st := store.NewJSONStore(*path)
		workspace, err := st.Load()
		if err != nil {
			return err
		}
		dispatch := service.NewDispatchService(workspace)
		if err := dispatch.CancelOrder(*orderID); err != nil {
			return err
		}
		return st.Save(dispatch.Workspace)
	case "report":
		fs := flag.NewFlagSet("report", flag.ContinueOnError)
		path := fs.String("file", "workspace.json", "workspace JSON path")
		nowText := fs.String("now", time.Now().UTC().Format(time.RFC3339), "report timestamp")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		now, err := time.Parse(time.RFC3339, *nowText)
		if err != nil {
			return fmt.Errorf("parse -now: %w", err)
		}
		workspace, err := store.NewJSONStore(*path).Load()
		if err != nil {
			return err
		}
		summary := report.BuildSummary(workspace, now)
		return json.NewEncoder(os.Stdout).Encode(summary)
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}
