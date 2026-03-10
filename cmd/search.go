package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cstsortan/prm/internal/render"
	"github.com/cstsortan/prm/internal/service"
)

func init() {
	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search across all entities",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			typeStr, _ := cmd.Flags().GetString("type")
			statusStr, _ := cmd.Flags().GetString("status")

			svc, err := getService()
			if err != nil {
				return err
			}

			results, err := svc.Search(args[0], service.ListFilter{
				Types:    service.ParseTypes(typeStr),
				Statuses: service.ParseStatuses(statusStr),
			})
			if err != nil {
				return err
			}

			if len(results) == 0 {
				fmt.Println(render.Dim.Render("No results found."))
				return nil
			}

			// Convert to ListResult for table rendering
			listResults := make([]service.ListResult, len(results))
			for i, r := range results {
				listResults[i] = service.ListResult{Entity: r.Entity, Path: r.Path}
			}

			fmt.Print(render.Table(listResults))
			return nil
		},
	}

	cmd.Flags().String("type", "", "Filter by type (comma-separated)")
	cmd.Flags().String("status", "", "Filter by status (comma-separated)")
	rootCmd.AddCommand(cmd)
}
