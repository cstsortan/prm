package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cstsortan/prm/internal/model"
	"github.com/cstsortan/prm/internal/tui"
)

func init() {
	subtaskCmd := &cobra.Command{
		Use:   "subtask",
		Short: "Manage sub-tasks",
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new sub-task under a story",
		RunE: func(cmd *cobra.Command, args []string) error {
			storyRef, _ := cmd.Flags().GetString("story")
			if storyRef == "" && !tui.IsInteractive() {
				return fmt.Errorf("--story is required")
			}
			return createEntity(cmd, model.EntitySubTask, storyRef)
		},
	}
	addCreateFlags(createCmd)
	createCmd.Flags().String("story", "", "Parent story (ID or slug, required)")
	subtaskCmd.AddCommand(createCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List sub-tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listEntities(cmd, model.EntitySubTask)
		},
	}
	addFilterFlags(listCmd)
	subtaskCmd.AddCommand(listCmd)

	subtaskCmd.AddCommand(&cobra.Command{
		Use:   "show <id-or-slug>",
		Short: "Show sub-task details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return showEntity(args[0])
		},
	})

	updateCmd := &cobra.Command{
		Use:   "update <id-or-slug>",
		Short: "Update a sub-task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateEntity(cmd, args[0])
		},
	}
	addUpdateFlags(updateCmd)
	subtaskCmd.AddCommand(updateCmd)

	subtaskCmd.AddCommand(&cobra.Command{
		Use:   "edit <id-or-slug>",
		Short: "Edit a sub-task's README.md in $EDITOR",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return editEntity(args[0])
		},
	})

	subtaskCmd.AddCommand(&cobra.Command{
		Use:   "delete <id-or-slug>",
		Short: "Delete a sub-task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteEntity(args[0])
		},
	})

	rootCmd.AddCommand(subtaskCmd)
}
