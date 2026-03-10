package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cstsortan/prm/internal/render"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "dashboard",
		Short: "Show project summary statistics",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := getService()
			if err != nil {
				return err
			}

			stats, cfg, err := svc.Dashboard()
			if err != nil {
				return err
			}

			fmt.Print(render.Dashboard(stats, cfg))
			return nil
		},
	})
}
