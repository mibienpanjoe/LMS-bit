package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mibienpanjoe/LMS-bit/internal/config"
	"github.com/mibienpanjoe/LMS-bit/internal/logging"
	"github.com/mibienpanjoe/LMS-bit/internal/ui/tui"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()
	logger := logging.New(cfg.LogLevel)

	program := tea.NewProgram(tui.NewModel(cfg, logger), tea.WithAltScreen())

	done := make(chan error, 1)
	go func() {
		_, err := program.Run()
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			fmt.Fprintf(os.Stderr, "application error: %v\n", err)
			os.Exit(1)
		}
	case <-ctx.Done():
		program.Quit()
		err := <-done
		if err != nil && !errors.Is(err, tea.ErrProgramKilled) {
			fmt.Fprintf(os.Stderr, "shutdown error: %v\n", err)
			os.Exit(1)
		}
	}
}
