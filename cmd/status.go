package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cstsortan/prm/internal/render"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "status <id-or-slug> <new-status>",
		Short: "Change an entity's status",
		Long:  "Valid statuses: backlog, todo, in-progress, review, done, cancelled, archived",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := getService()
			if err != nil {
				return err
			}

			entity, err := svc.SetStatus(args[0], args[1])
			if err != nil {
				return err
			}

			fmt.Printf("%s %s -> %s\n",
				render.Success.Render("OK"),
				entity.Title,
				render.StatusStyle(entity.Status),
			)
			return nil
		},
	})
}
