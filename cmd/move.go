package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cstsortan/prm/internal/render"
)

func init() {
	cmd := &cobra.Command{
		Use:   "move <id-or-slug>",
		Short: "Move an entity to a different parent",
		Long:  "Move a story to a different epic, a sub-task to a different story, or promote a sub-task to a standalone task with --standalone.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			to, _ := cmd.Flags().GetString("to")
			standalone, _ := cmd.Flags().GetBool("standalone")

			if to == "" && !standalone {
				return fmt.Errorf("either --to <parent> or --standalone is required")
			}
			if to != "" && standalone {
				return fmt.Errorf("--to and --standalone are mutually exclusive")
			}

			svc, err := getService()
			if err != nil {
				return err
			}

			entity, err := svc.MoveEntity(args[0], to, standalone)
			if err != nil {
				return err
			}

			if standalone {
				fmt.Printf("%s Promoted %s to standalone task\n",
					render.Success.Render("OK"),
					entity.Title,
				)
			} else {
				fmt.Printf("%s Moved %s: %s\n",
					render.Success.Render("OK"),
					render.TypeStyle(entity.Type),
					entity.Title,
				)
			}
			return nil
		},
	}

	cmd.Flags().String("to", "", "New parent entity (ID or slug)")
	cmd.Flags().Bool("standalone", false, "Promote sub-task to standalone task")
	rootCmd.AddCommand(cmd)
}
