package web

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/cstsortan/prm/internal/service"
)

// Server is the PRM web UI HTTP server.
type Server struct {
	svc  *service.Service
	port int
	srv  *http.Server
}

// NewServer creates a new web server wrapping the given service.
func NewServer(svc *service.Service, port int) *Server {
	s := &Server{svc: svc, port: port}
	mux := http.NewServeMux()
	s.registerRoutes(mux)
	s.srv = &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	return s
}

// Start starts the server and optionally opens the browser.
func (s *Server) Start(openBrowser bool) error {
	ln, err := net.Listen("tcp", s.srv.Addr)
	if err != nil {
		return fmt.Errorf("listening on port %d: %w", s.port, err)
	}

	url := fmt.Sprintf("http://localhost:%d", s.port)
	fmt.Printf("PRM Web UI running at %s\n", url)
	fmt.Println("Press Ctrl+C to stop")

	if openBrowser {
		openURL(url)
	}

	if err := s.srv.Serve(ln); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func openURL(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	}
	if cmd != nil {
		cmd.Start()
	}
}
