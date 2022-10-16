package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bankparse",
		Short: "",
		Long:  ``,
	}
	cmd.PersistentFlags().StringVar(&Config.ConfigFile, "config", "", "Config file (default: ~/.config/bankparse/config.toml)")

	cmd.AddCommand(NewParseCmd())
	return cmd
}

func main() {
	if err := NewRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
