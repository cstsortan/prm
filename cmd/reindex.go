package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cstsortan/prm/internal/render"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "reindex",
		Short: "Rebuild index.json from the file tree",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := getService()
			if err != nil {
				return err
			}

			idx, err := svc.Store.RebuildIndex()
			if err != nil {
				return err
			}

			fmt.Printf("%s Reindexed %d entities\n",
				render.Success.Render("OK"),
				len(idx.Entries),
			)
			return nil
		},
	})
}
