package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cstsortan/prm/internal/render"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "tree [epic-id-or-slug]",
		Short: "Show hierarchy tree",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := getService()
			if err != nil {
				return err
			}

			ref := ""
			if len(args) > 0 {
				ref = args[0]
			}

			nodes, err := svc.Tree(ref)
			if err != nil {
				return err
			}

			fmt.Print(render.Tree(nodes))
			return nil
		},
	})
}
