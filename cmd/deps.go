package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cstsortan/prm/internal/render"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "deps <id-or-slug>",
		Short: "Show dependencies for an entity",
		Long:  "Shows what this entity depends on and what depends on it.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := getService()
			if err != nil {
				return err
			}

			deps, err := svc.GetDependencies(args[0])
			if err != nil {
				return err
			}

			var b strings.Builder
			b.WriteString(render.Title.Render(deps.Entity.Title))
			b.WriteString("\n\n")

			if len(deps.DependsOn) == 0 && len(deps.DependedBy) == 0 {
				b.WriteString(render.Dim.Render("  No dependencies."))
				b.WriteString("\n")
			}

			if len(deps.DependsOn) > 0 {
				b.WriteString(render.Bold.Render("  Depends on:"))
				b.WriteString("\n")
				for _, d := range deps.DependsOn {
					b.WriteString(fmt.Sprintf("    %s %s %s\n",
						render.StatusStyle(d.Status),
						d.Title,
						render.IDStyle.Render("("+d.ShortID+")"),
					))
				}
			}

			if len(deps.DependedBy) > 0 {
				b.WriteString(render.Bold.Render("  Depended on by:"))
				b.WriteString("\n")
				for _, d := range deps.DependedBy {
					b.WriteString(fmt.Sprintf("    %s %s %s\n",
						render.StatusStyle(d.Status),
						d.Title,
						render.IDStyle.Render("("+d.ShortID+")"),
					))
				}
			}

			fmt.Print(b.String())
			return nil
		},
	})
}
