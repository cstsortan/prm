package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/cstsortan/prm/internal/model"
	"github.com/cstsortan/prm/internal/render"
	"github.com/cstsortan/prm/internal/service"
	"github.com/cstsortan/prm/internal/tui"
)

func init() {
	epicCmd := &cobra.Command{
		Use:   "epic",
		Short: "Manage epics",
	}

	// create
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new epic",
		RunE: func(cmd *cobra.Command, args []string) error {
			return createEntity(cmd, model.EntityEpic, "")
		},
	}
	addCreateFlags(createCmd)
	epicCmd.AddCommand(createCmd)

	// list
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List epics",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listEntities(cmd, model.EntityEpic)
		},
	}
	addFilterFlags(listCmd)
	epicCmd.AddCommand(listCmd)

	// show
	epicCmd.AddCommand(&cobra.Command{
		Use:   "show <id-or-slug>",
		Short: "Show epic details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return showEntity(args[0])
		},
	})

	// update
	updateCmd := &cobra.Command{
		Use:   "update <id-or-slug>",
		Short: "Update an epic",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateEntity(cmd, args[0])
		},
	}
	addUpdateFlags(updateCmd)
	epicCmd.AddCommand(updateCmd)

	// edit
	epicCmd.AddCommand(&cobra.Command{
		Use:   "edit <id-or-slug>",
		Short: "Edit an epic's README.md in $EDITOR",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return editEntity(args[0])
		},
	})

	// delete
	epicCmd.AddCommand(&cobra.Command{
		Use:   "delete <id-or-slug>",
		Short: "Delete an epic",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteEntity(args[0])
		},
	})

	rootCmd.AddCommand(epicCmd)
}

// Shared helpers for all entity commands

func addCreateFlags(cmd *cobra.Command) {
	cmd.Flags().String("title", "", "Title (required)")
	cmd.Flags().String("description", "", "Short description")
	cmd.Flags().String("body", "", "README.md content (detailed description)")
	cmd.Flags().String("priority", "medium", "Priority (low, medium, high, critical)")
	cmd.Flags().String("status", "backlog", "Initial status")
	cmd.Flags().String("tags", "", "Comma-separated tags")
	cmd.Flags().String("due", "", "Due date (YYYY-MM-DD)")
	cmd.Flags().String("depends", "", "Comma-separated IDs/slugs this item depends on")
}

func addFilterFlags(cmd *cobra.Command) {
	cmd.Flags().String("status", "", "Filter by status (comma-separated)")
	cmd.Flags().String("priority", "", "Filter by priority (comma-separated)")
	cmd.Flags().String("tag", "", "Filter by tag (comma-separated)")
}

func addUpdateFlags(cmd *cobra.Command) {
	cmd.Flags().String("title", "", "New title")
	cmd.Flags().String("description", "", "New description")
	cmd.Flags().String("body", "", "New README.md content (detailed description)")
	cmd.Flags().String("priority", "", "New priority")
	cmd.Flags().String("tags", "", "New tags (comma-separated)")
	cmd.Flags().String("due", "", "New due date (YYYY-MM-DD)")
	cmd.Flags().Bool("clear-due", false, "Remove due date")
	cmd.Flags().String("severity", "", "New severity (bugs only)")
	cmd.Flags().String("depends", "", "Dependencies (comma-separated IDs/slugs)")
	cmd.Flags().Bool("clear-depends", false, "Remove all dependencies")
}

func createEntity(cmd *cobra.Command, entityType model.EntityType, parentFlag string) error {
	title := strings.TrimSpace(mustGetString(cmd, "title"))

	// Interactive wizard when title is missing and terminal is available
	if title == "" && tui.IsInteractive() {
		return createEntityInteractive(entityType, parentFlag)
	}

	if title == "" {
		return fmt.Errorf("--title is required")
	}

	// If parent is needed but missing, prompt interactively
	if parentFlag == "" && tui.IsInteractive() {
		if entityType == model.EntityStory || entityType == model.EntitySubTask {
			ref, err := promptParentSelect(entityType)
			if err != nil {
				return err
			}
			parentFlag = ref
		}
	}

	description, _ := cmd.Flags().GetString("description")
	body, _ := cmd.Flags().GetString("body")
	priorityStr, _ := cmd.Flags().GetString("priority")
	statusStr, _ := cmd.Flags().GetString("status")
	tagsStr, _ := cmd.Flags().GetString("tags")
	dueStr, _ := cmd.Flags().GetString("due")
	dependsStr, _ := cmd.Flags().GetString("depends")

	priority, ok := model.ParsePriority(priorityStr)
	if !ok {
		return fmt.Errorf("invalid priority: %s", priorityStr)
	}
	status, ok := model.ParseStatus(statusStr)
	if !ok {
		return fmt.Errorf("invalid status: %s", statusStr)
	}

	var dueDate *time.Time
	if dueStr != "" {
		t, err := time.Parse("2006-01-02", dueStr)
		if err != nil {
			return fmt.Errorf("invalid due date format (use YYYY-MM-DD): %w", err)
		}
		dueDate = &t
	}

	// Bug-specific flags
	var severity model.Severity
	if entityType == model.EntityBug {
		sevStr, _ := cmd.Flags().GetString("severity")
		if sevStr != "" {
			sv, ok := model.ParseSeverity(sevStr)
			if !ok {
				return fmt.Errorf("invalid severity: %s", sevStr)
			}
			severity = sv
		}
	}

	svc, err := getService()
	if err != nil {
		return err
	}

	entity, err := svc.CreateEntity(service.CreateEntityOpts{
		Type:         entityType,
		Title:        title,
		Description:  description,
		Body:         body,
		Priority:     priority,
		Status:       status,
		Tags:         service.ParseTags(tagsStr),
		DueDate:      dueDate,
		ParentID:     parentFlag,
		Severity:     severity,
		Dependencies: service.ParseTags(dependsStr),
	})
	if err != nil {
		return err
	}

	fmt.Printf("%s Created %s: %s %s\n",
		render.Success.Render("OK"),
		render.TypeStyle(entity.Type),
		entity.Title,
		render.IDStyle.Render("("+entity.ShortID()+")"),
	)
	return nil
}

func createEntityInteractive(entityType model.EntityType, parentFlag string) error {
	svc, err := getService()
	if err != nil {
		return err
	}

	// Build wizard steps
	steps := []tui.WizardStep{
		{Label: "Title:", Kind: tui.StepText, Required: true, Placeholder: "Enter a title..."},
		{Label: "Description:", Kind: tui.StepText, Placeholder: "Short description (optional)"},
		{Label: "Priority:", Kind: tui.StepSelect, Options: prioritySelectOptions(), DefaultIdx: 1},
		{Label: "Status:", Kind: tui.StepSelect, Options: statusSelectOptions(), DefaultIdx: 0},
		{Label: "Tags:", Kind: tui.StepText, Placeholder: "backend, frontend, ..."},
		{Label: "Due date:", Kind: tui.StepText, Placeholder: "YYYY-MM-DD"},
	}

	// Bug severity
	if entityType == model.EntityBug {
		steps = append(steps, tui.WizardStep{
			Label: "Severity:", Kind: tui.StepSelect, Options: severitySelectOptions(), DefaultIdx: 1,
		})
	}

	// Parent entity selection
	if parentFlag == "" && (entityType == model.EntityStory || entityType == model.EntitySubTask) {
		choices, parentLabel, err := buildParentChoices(svc, entityType)
		if err != nil {
			return err
		}
		steps = append(steps, tui.WizardStep{
			Label: parentLabel, Kind: tui.StepFilter, Choices: choices,
		})
	}

	result, err := tui.RunWizard(steps)
	if err != nil {
		return err
	}
	if result.Cancelled {
		fmt.Println("Cancelled.")
		return nil
	}

	// Parse wizard values
	title := result.Values[0]
	description := result.Values[1]
	priority, _ := model.ParsePriority(result.Values[2])
	status, _ := model.ParseStatus(result.Values[3])
	tags := service.ParseTags(result.Values[4])

	var dueDate *time.Time
	if result.Values[5] != "" {
		t, err := time.Parse("2006-01-02", result.Values[5])
		if err != nil {
			return fmt.Errorf("invalid due date format (use YYYY-MM-DD): %w", err)
		}
		dueDate = &t
	}

	var severity model.Severity
	parentRef := parentFlag
	idx := 6

	if entityType == model.EntityBug {
		severity, _ = model.ParseSeverity(result.Values[idx])
		idx++
	}

	if parentFlag == "" && idx < len(result.Values) {
		parentRef = result.Values[idx]
	}

	entity, err := svc.CreateEntity(service.CreateEntityOpts{
		Type:        entityType,
		Title:       title,
		Description: description,
		Priority:    priority,
		Status:      status,
		Tags:        tags,
		DueDate:     dueDate,
		ParentID:    parentRef,
		Severity:    severity,
	})
	if err != nil {
		return err
	}

	fmt.Printf("%s Created %s: %s %s\n",
		render.Success.Render("OK"),
		render.TypeStyle(entity.Type),
		entity.Title,
		render.IDStyle.Render("("+entity.ShortID()+")"),
	)
	return nil
}

func promptParentSelect(entityType model.EntityType) (string, error) {
	svc, err := getService()
	if err != nil {
		return "", err
	}

	choices, label, err := buildParentChoices(svc, entityType)
	if err != nil {
		return "", err
	}

	return tui.PromptEntitySelect(label, choices)
}

func buildParentChoices(svc *service.Service, entityType model.EntityType) ([]tui.EntityChoice, string, error) {
	var parentType model.EntityType
	var label string

	switch entityType {
	case model.EntityStory:
		parentType = model.EntityEpic
		label = "Select parent epic:"
	case model.EntitySubTask:
		parentType = model.EntityStory
		label = "Select parent story:"
	default:
		return nil, "", fmt.Errorf("entity type %s has no parent", entityType)
	}

	parents, err := svc.List(service.ListFilter{Types: []model.EntityType{parentType}})
	if err != nil {
		return nil, "", err
	}
	if len(parents) == 0 {
		return nil, "", fmt.Errorf("no %ss found; create one first", parentType)
	}

	choices := make([]tui.EntityChoice, len(parents))
	for i, p := range parents {
		choices[i] = tui.EntityChoice{
			Label: fmt.Sprintf("%s (%s)", p.Entity.Title, p.Entity.ShortID()),
			Value: p.Entity.Slug,
		}
	}
	return choices, label, nil
}

func prioritySelectOptions() []tui.SelectOption {
	return []tui.SelectOption{
		{Label: "low", Value: "low"},
		{Label: "medium", Value: "medium"},
		{Label: "high", Value: "high"},
		{Label: "critical", Value: "critical"},
	}
}

func statusSelectOptions() []tui.SelectOption {
	return []tui.SelectOption{
		{Label: "backlog", Value: "backlog"},
		{Label: "todo", Value: "todo"},
		{Label: "in-progress", Value: "in-progress"},
		{Label: "review", Value: "review"},
		{Label: "done", Value: "done"},
		{Label: "cancelled", Value: "cancelled"},
	}
}

func severitySelectOptions() []tui.SelectOption {
	return []tui.SelectOption{
		{Label: "cosmetic", Value: "cosmetic"},
		{Label: "minor", Value: "minor"},
		{Label: "major", Value: "major"},
		{Label: "blocker", Value: "blocker"},
	}
}

func mustGetString(cmd *cobra.Command, name string) string {
	val, _ := cmd.Flags().GetString(name)
	return val
}

func listEntities(cmd *cobra.Command, entityType model.EntityType) error {
	statusStr, _ := cmd.Flags().GetString("status")
	priorityStr, _ := cmd.Flags().GetString("priority")
	tagStr, _ := cmd.Flags().GetString("tag")

	statuses := service.ParseStatuses(statusStr)
	if err := service.ValidateStatuses(statuses); err != nil {
		return err
	}
	priorities := service.ParsePriorities(priorityStr)
	if err := service.ValidatePriorities(priorities); err != nil {
		return err
	}

	svc, err := getService()
	if err != nil {
		return err
	}

	results, err := svc.List(service.ListFilter{
		Types:      []model.EntityType{entityType},
		Statuses:   statuses,
		Priorities: priorities,
		Tags:       service.ParseTags(tagStr),
	})
	if err != nil {
		return err
	}

	fmt.Print(render.Table(results))
	return nil
}

func showEntity(ref string) error {
	svc, err := getService()
	if err != nil {
		return err
	}

	entity, _, readme, depMap, err := svc.ShowEntity(ref)
	if err != nil {
		return err
	}

	fmt.Print(render.EntityDetail(entity, readme, depMap))
	return nil
}

func updateEntity(cmd *cobra.Command, ref string) error {
	title, _ := cmd.Flags().GetString("title")
	description, _ := cmd.Flags().GetString("description")
	body, _ := cmd.Flags().GetString("body")
	priorityStr, _ := cmd.Flags().GetString("priority")
	tagsStr, _ := cmd.Flags().GetString("tags")
	dueStr, _ := cmd.Flags().GetString("due")
	clearDue, _ := cmd.Flags().GetBool("clear-due")
	severityStr, _ := cmd.Flags().GetString("severity")
	dependsStr, _ := cmd.Flags().GetString("depends")
	clearDepends, _ := cmd.Flags().GetBool("clear-depends")

	var dueDate *time.Time
	if dueStr != "" {
		t, err := time.Parse("2006-01-02", dueStr)
		if err != nil {
			return fmt.Errorf("invalid due date format (use YYYY-MM-DD): %w", err)
		}
		dueDate = &t
	}

	var tags []string
	if cmd.Flags().Changed("tags") {
		tags = service.ParseTags(tagsStr)
	}

	var deps []string
	if cmd.Flags().Changed("depends") {
		deps = service.ParseTags(dependsStr)
	}

	svc, err := getService()
	if err != nil {
		return err
	}

	entity, err := svc.UpdateEntity(ref, service.UpdateEntityOpts{
		Title:        title,
		Description:  description,
		Body:         body,
		Priority:     priorityStr,
		Tags:         tags,
		DueDate:      dueDate,
		ClearDue:     clearDue,
		Severity:     severityStr,
		Dependencies: deps,
		ClearDepends: clearDepends,
	})
	if err != nil {
		return err
	}

	fmt.Printf("%s Updated %s: %s\n",
		render.Success.Render("OK"),
		render.TypeStyle(entity.Type),
		entity.Title,
	)
	return nil
}

func editEntity(ref string) error {
	svc, err := getService()
	if err != nil {
		return err
	}

	_, dir, _, _, err := svc.ShowEntity(ref)
	if err != nil {
		return err
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	readmePath := dir + "/README.md"
	cmd := exec.Command(editor, readmePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("editor exited with error: %w", err)
	}

	fmt.Printf("%s Edited README.md\n", render.Success.Render("OK"))
	return nil
}

func deleteEntity(ref string) error {
	if tui.IsInteractive() {
		confirmed, err := tui.Confirm("Delete this entity and all its children?")
		if err != nil {
			return err
		}
		if !confirmed {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	svc, err := getService()
	if err != nil {
		return err
	}

	if err := svc.DeleteEntity(ref); err != nil {
		return err
	}

	fmt.Printf("%s Deleted\n", render.Success.Render("OK"))
	return nil
}
