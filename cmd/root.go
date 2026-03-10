package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/cstsortan/prm/internal/service"
	"github.com/cstsortan/prm/internal/store"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:     "prm",
	Short:   "PRM - Project Resource Manager",
	Long:    "A CLI project management tool that stores data as .json and .md files.",
	Version: version,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// getService creates a Service instance using the current working directory.
// Returns an error if .prm/ is not initialized.
func getService() (*service.Service, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("getting working directory: %w", err)
	}
	s := store.New(dir)
	if !s.Exists() {
		return nil, fmt.Errorf("prm not initialized (run 'prm init' first)")
	}
	return service.New(s), nil
}

// getServiceForInit creates a Service instance without checking if .prm/ exists.
func getServiceForInit() (*service.Service, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("getting working directory: %w", err)
	}
	s := store.New(dir)
	return service.New(s), nil
}
