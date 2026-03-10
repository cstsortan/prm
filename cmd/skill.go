package cmd

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/cstsortan/prm/internal/render"
)

//go:embed skill_prm.md
var skillContent string

func init() {
	cmd := &cobra.Command{
		Use:   "install-skill",
		Short: "Install the PRM Claude Code skill into this project",
		Long:  "Writes the PRM skill file to .claude/commands/prm.md so Claude Code can manage work items.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			force, _ := cmd.Flags().GetBool("force")

			dir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("getting working directory: %w", err)
			}

			skillDir := filepath.Join(dir, ".claude", "commands")
			skillPath := filepath.Join(skillDir, "prm.md")

			if !force {
				if _, err := os.Stat(skillPath); err == nil {
					return fmt.Errorf("skill file already exists at %s (use --force to overwrite)", skillPath)
				}
			}

			if err := os.MkdirAll(skillDir, 0755); err != nil {
				return fmt.Errorf("creating directory %s: %w", skillDir, err)
			}

			if err := os.WriteFile(skillPath, []byte(skillContent), 0644); err != nil {
				return fmt.Errorf("writing skill file: %w", err)
			}

			fmt.Println(render.Success.Render("Installed PRM skill to .claude/commands/prm.md"))
			return nil
		},
	}

	cmd.Flags().Bool("force", false, "Overwrite existing skill file")
	rootCmd.AddCommand(cmd)
}
