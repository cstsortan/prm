package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/cstsortan/prm/internal/render"
)

func init() {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export entities to CSV or JSON",
		RunE: func(cmd *cobra.Command, args []string) error {
			format, _ := cmd.Flags().GetString("format")
			output, _ := cmd.Flags().GetString("output")

			svc, err := getService()
			if err != nil {
				return err
			}

			var w *os.File
			if output != "" {
				w, err = os.Create(output)
				if err != nil {
					return fmt.Errorf("creating output file: %w", err)
				}
				defer w.Close()
			} else {
				w = os.Stdout
			}

			switch format {
			case "csv":
				if err := svc.ExportCSV(w); err != nil {
					return err
				}
			case "json":
				if err := svc.ExportJSON(w); err != nil {
					return err
				}
			default:
				return fmt.Errorf("invalid format: %s (use csv or json)", format)
			}

			if output != "" {
				fmt.Printf("%s Exported to %s\n", render.Success.Render("OK"), output)
			}
			return nil
		},
	}

	cmd.Flags().String("format", "csv", "Export format (csv, json)")
	cmd.Flags().String("output", "", "Output file (defaults to stdout)")
	rootCmd.AddCommand(cmd)
}
