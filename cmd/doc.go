package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cstsortan/prm/internal/render"
	"github.com/cstsortan/prm/internal/store"
	"github.com/cstsortan/prm/internal/tui"
)

func init() {
	docCmd := &cobra.Command{
		Use:   "doc",
		Short: "Manage documents",
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new document",
		RunE: func(cmd *cobra.Command, args []string) error {
			title, _ := cmd.Flags().GetString("title")
			body, _ := cmd.Flags().GetString("body")
			if title == "" && tui.IsInteractive() {
				t, err := tui.PromptText("Document title:", "Enter a title...", true)
				if err != nil {
					return err
				}
				title = t
			}
			if title == "" {
				return fmt.Errorf("--title is required")
			}

			svc, err := getService()
			if err != nil {
				return err
			}

			docsDir := filepath.Join(svc.Store.Root(), "docs")
			slug := store.UniqueSlug(docsDir, store.GenerateSlug(title))
			filename := slug + ".md"
			path := filepath.Join(docsDir, filename)

			var content string
			if body != "" {
				content = body + "\n"
			} else {
				content = fmt.Sprintf("# %s\n\n", title)
			}
			if err := store.WriteFileAtomic(path, []byte(content)); err != nil {
				return fmt.Errorf("creating doc: %w", err)
			}

			fmt.Printf("%s Created doc: %s\n", render.Success.Render("OK"), filename)
			return nil
		},
	}
	createCmd.Flags().String("title", "", "Document title (required)")
	createCmd.Flags().String("body", "", "Document content (markdown)")
	docCmd.AddCommand(createCmd)

	docCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List documents",
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := getService()
			if err != nil {
				return err
			}

			docsDir := filepath.Join(svc.Store.Root(), "docs")
			entries, err := os.ReadDir(docsDir)
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Println(render.Dim.Render("No documents found."))
					return nil
				}
				return err
			}

			found := false
			for _, e := range entries {
				if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
					fmt.Printf("  %s\n", e.Name())
					found = true
				}
			}
			if !found {
				fmt.Println(render.Dim.Render("No documents found."))
			}
			return nil
		},
	})

	docCmd.AddCommand(&cobra.Command{
		Use:   "show <slug>",
		Short: "Show a document's contents",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			slug := args[0]
			if strings.ContainsAny(slug, "/\\..") {
				return fmt.Errorf("invalid document slug: %q", slug)
			}

			svc, err := getService()
			if err != nil {
				return err
			}

			if !strings.HasSuffix(slug, ".md") {
				slug += ".md"
			}
			path := filepath.Join(svc.Store.Root(), "docs", slug)

			data, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("reading doc: %w", err)
			}

			fmt.Print(string(data))
			return nil
		},
	})

	rootCmd.AddCommand(docCmd)
}
