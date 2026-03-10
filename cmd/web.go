package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/cstsortan/prm/internal/web"
)

func init() {
	rootCmd.AddCommand(webCmd)
	webCmd.Flags().IntP("port", "p", 3141, "port to listen on")
	webCmd.Flags().Bool("no-open", false, "don't open browser automatically")
}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Launch the PRM web UI",
	Long:  "Starts a local web server serving the PRM dashboard and API.",
	RunE: func(cmd *cobra.Command, args []string) error {
		svc, err := getService()
		if err != nil {
			return err
		}

		port, _ := cmd.Flags().GetInt("port")
		noOpen, _ := cmd.Flags().GetBool("no-open")

		server := web.NewServer(svc, port)

		// Graceful shutdown on Ctrl+C
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		errCh := make(chan error, 1)
		go func() {
			errCh <- server.Start(!noOpen)
		}()

		select {
		case err := <-errCh:
			return err
		case <-ctx.Done():
			fmt.Println("\nShutting down...")
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			return server.Shutdown(shutdownCtx)
		}
	},
}
