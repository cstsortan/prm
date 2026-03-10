package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cstsortan/prm/internal/render"
	"github.com/cstsortan/prm/internal/service"
)

func init() {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all entities with optional filters",
		RunE: func(cmd *cobra.Command, args []string) error {
			typeStr, _ := cmd.Flags().GetString("type")
			statusStr, _ := cmd.Flags().GetString("status")
			priorityStr, _ := cmd.Flags().GetString("priority")
			tagStr, _ := cmd.Flags().GetString("tag")
			sortStr, _ := cmd.Flags().GetString("sort")
			sortDesc, _ := cmd.Flags().GetBool("desc")
			archived, _ := cmd.Flags().GetBool("archived")

			statuses := service.ParseStatuses(statusStr)
			if err := service.ValidateStatuses(statuses); err != nil {
				return err
			}
			priorities := service.ParsePriorities(priorityStr)
			if err := service.ValidatePriorities(priorities); err != nil {
				return err
			}

			var sortField service.SortField
			if sortStr != "" {
				sf, ok := service.ParseSortField(sortStr)
				if !ok {
					return fmt.Errorf("invalid sort field: %q (valid: priority, status, title, created, updated, type)", sortStr)
				}
				sortField = sf
			}

			svc, err := getService()
			if err != nil {
				return err
			}

			results, err := svc.List(service.ListFilter{
				Types:           service.ParseTypes(typeStr),
				Statuses:        statuses,
				Priorities:      priorities,
				Tags:            service.ParseTags(tagStr),
				Sort:            sortField,
				SortDesc:        sortDesc,
				IncludeArchived: archived,
			})
			if err != nil {
				return err
			}

			fmt.Print(render.Table(results))
			return nil
		},
	}

	cmd.Flags().String("type", "", "Filter by type (comma-separated: epic,story,sub-task,task,bug)")
	cmd.Flags().String("status", "", "Filter by status (comma-separated)")
	cmd.Flags().String("priority", "", "Filter by priority (comma-separated)")
	cmd.Flags().String("tag", "", "Filter by tag (comma-separated)")
	cmd.Flags().String("sort", "", "Sort by field (priority, status, title, created, updated, type)")
	cmd.Flags().Bool("desc", false, "Reverse sort order")
	cmd.Flags().Bool("archived", false, "Include archived items")
	rootCmd.AddCommand(cmd)
}
