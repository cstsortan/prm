package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cstsortan/prm/internal/render"
)

func init() {
	cmd := &cobra.Command{
		Use:   "comment <id-or-slug>",
		Short: "Add a comment to an entity",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			text, _ := cmd.Flags().GetString("text")
			if text == "" {
				return fmt.Errorf("--text is required")
			}
			author, _ := cmd.Flags().GetString("author")

			svc, err := getService()
			if err != nil {
				return err
			}

			entity, err := svc.AddComment(args[0], author, text)
			if err != nil {
				return err
			}

			fmt.Printf("%s Comment added to %s\n",
				render.Success.Render("OK"),
				entity.Title,
			)
			return nil
		},
	}

	cmd.Flags().String("text", "", "Comment text (required)")
	cmd.Flags().String("author", "user", "Comment author")
	rootCmd.AddCommand(cmd)
}
