package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cstsortan/prm/internal/render"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "archive <id-or-slug>",
		Short: "Archive an entity",
		Long:  "Sets the entity's status to archived. Archived items are hidden from list and dashboard by default. Use 'prm list --archived' to see them.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := getService()
			if err != nil {
				return err
			}

			entity, err := svc.SetStatus(args[0], "archived")
			if err != nil {
				return err
			}

			fmt.Printf("%s Archived: %s\n",
				render.Success.Render("OK"),
				entity.Title,
			)
			return nil
		},
	})
}
