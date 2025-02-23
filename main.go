package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/katexochen/ghh/internal/cmd"

	"github.com/spf13/cobra"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
}

func run() error {
	cobra.EnableCommandSorting = false
	rootCmd := newRootCmd()
	ctx, cancel := signalContext(context.Background(), os.Interrupt)
	defer cancel()
	return rootCmd.ExecuteContext(ctx)
}

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "ghh",
		Short: "GitHub Helper CLI",
	}

	rootCmd.SetOut(os.Stdout)
	rootCmd.AddCommand(
		cmd.NewDeleteAllRunsCmd(),
		cmd.NewCreateProjectIssueCmd(),
		cmd.NewSetAuthCmd(),
		cmd.NewListBranchesCmd(),
	)
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")

	return rootCmd
}

func signalContext(ctx context.Context, sig os.Signal) (context.Context, context.CancelFunc) {
	sigCtx, stop := signal.NotifyContext(ctx, sig)
	done := make(chan struct{}, 1)
	stopDone := make(chan struct{}, 1)

	go func() {
		defer func() { stopDone <- struct{}{} }()
		defer stop()
		select {
		case <-sigCtx.Done():
			fmt.Println(" Signal caught. Press ctrl+c again to terminate the program immediately.")
		case <-done:
		}
	}()

	cancelFunc := func() {
		done <- struct{}{}
		<-stopDone
	}

	return sigCtx, cancelFunc
}
