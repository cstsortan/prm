package cmd

import (
	"github.com/spf13/cobra"

	"github.com/cstsortan/prm/internal/model"
)

func init() {
	bugCmd := &cobra.Command{
		Use:   "bug",
		Short: "Manage bugs",
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new bug",
		RunE: func(cmd *cobra.Command, args []string) error {
			return createEntity(cmd, model.EntityBug, "")
		},
	}
	addCreateFlags(createCmd)
	createCmd.Flags().String("severity", "", "Severity (cosmetic, minor, major, blocker)")
	bugCmd.AddCommand(createCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List bugs",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listEntities(cmd, model.EntityBug)
		},
	}
	addFilterFlags(listCmd)
	bugCmd.AddCommand(listCmd)

	bugCmd.AddCommand(&cobra.Command{
		Use:   "show <id-or-slug>",
		Short: "Show bug details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return showEntity(args[0])
		},
	})

	updateCmd := &cobra.Command{
		Use:   "update <id-or-slug>",
		Short: "Update a bug",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateEntity(cmd, args[0])
		},
	}
	addUpdateFlags(updateCmd)
	bugCmd.AddCommand(updateCmd)

	bugCmd.AddCommand(&cobra.Command{
		Use:   "edit <id-or-slug>",
		Short: "Edit a bug's README.md in $EDITOR",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return editEntity(args[0])
		},
	})

	bugCmd.AddCommand(&cobra.Command{
		Use:   "delete <id-or-slug>",
		Short: "Delete a bug",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteEntity(args[0])
		},
	})

	rootCmd.AddCommand(bugCmd)
}
