package cmd

import (
	"github.com/spf13/cobra"

	"github.com/cstsortan/prm/internal/model"
)

func init() {
	taskCmd := &cobra.Command{
		Use:   "task",
		Short: "Manage standalone tasks",
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new standalone task",
		RunE: func(cmd *cobra.Command, args []string) error {
			return createEntity(cmd, model.EntityTask, "")
		},
	}
	addCreateFlags(createCmd)
	taskCmd.AddCommand(createCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List standalone tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listEntities(cmd, model.EntityTask)
		},
	}
	addFilterFlags(listCmd)
	taskCmd.AddCommand(listCmd)

	taskCmd.AddCommand(&cobra.Command{
		Use:   "show <id-or-slug>",
		Short: "Show task details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return showEntity(args[0])
		},
	})

	updateCmd := &cobra.Command{
		Use:   "update <id-or-slug>",
		Short: "Update a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateEntity(cmd, args[0])
		},
	}
	addUpdateFlags(updateCmd)
	taskCmd.AddCommand(updateCmd)

	taskCmd.AddCommand(&cobra.Command{
		Use:   "edit <id-or-slug>",
		Short: "Edit a task's README.md in $EDITOR",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return editEntity(args[0])
		},
	})

	taskCmd.AddCommand(&cobra.Command{
		Use:   "delete <id-or-slug>",
		Short: "Delete a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteEntity(args[0])
		},
	})

	rootCmd.AddCommand(taskCmd)
}
