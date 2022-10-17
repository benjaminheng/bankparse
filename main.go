package main

import (
	"os"

	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bankparse",
		Short: "",
		Long:  ``,
	}
	cmd.SilenceUsage = true

	cmd.AddCommand(NewParseCmd())
	return cmd
}

func main() {
	if err := NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
