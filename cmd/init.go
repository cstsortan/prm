package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/cstsortan/prm/internal/render"
)

func init() {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new PRM project in the current directory",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			if name == "" {
				abs, err := filepath.Abs(".")
				if err != nil {
					return fmt.Errorf("getting working directory: %w", err)
				}
				name = filepath.Base(abs)
			}

			svc, err := getServiceForInit()
			if err != nil {
				return err
			}

			if err := svc.Init(name); err != nil {
				return err
			}

			fmt.Println(render.Success.Render(fmt.Sprintf("Initialized PRM project: %s", name)))
			return nil
		},
	}

	cmd.Flags().String("name", "", "Project name (defaults to directory name)")
	rootCmd.AddCommand(cmd)
}
