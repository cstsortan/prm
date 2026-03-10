package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cstsortan/prm/internal/model"
	"github.com/cstsortan/prm/internal/tui"
)

func init() {
	storyCmd := &cobra.Command{
		Use:   "story",
		Short: "Manage stories",
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new story under an epic",
		RunE: func(cmd *cobra.Command, args []string) error {
			epicRef, _ := cmd.Flags().GetString("epic")
			if epicRef == "" && !tui.IsInteractive() {
				return fmt.Errorf("--epic is required")
			}
			return createEntity(cmd, model.EntityStory, epicRef)
		},
	}
	addCreateFlags(createCmd)
	createCmd.Flags().String("epic", "", "Parent epic (ID or slug, required)")
	storyCmd.AddCommand(createCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List stories",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listEntities(cmd, model.EntityStory)
		},
	}
	addFilterFlags(listCmd)
	storyCmd.AddCommand(listCmd)

	storyCmd.AddCommand(&cobra.Command{
		Use:   "show <id-or-slug>",
		Short: "Show story details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return showEntity(args[0])
		},
	})

	updateCmd := &cobra.Command{
		Use:   "update <id-or-slug>",
		Short: "Update a story",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateEntity(cmd, args[0])
		},
	}
	addUpdateFlags(updateCmd)
	storyCmd.AddCommand(updateCmd)

	storyCmd.AddCommand(&cobra.Command{
		Use:   "edit <id-or-slug>",
		Short: "Edit a story's README.md in $EDITOR",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return editEntity(args[0])
		},
	})

	storyCmd.AddCommand(&cobra.Command{
		Use:   "delete <id-or-slug>",
		Short: "Delete a story",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteEntity(args[0])
		},
	})

	rootCmd.AddCommand(storyCmd)
}
