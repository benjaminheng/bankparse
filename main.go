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
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := Config.Load(Config.ConfigFile)
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.SilenceUsage = true
	cmd.PersistentFlags().StringVar(&Config.ConfigFile, "config", "", "Config file (default: ~/.config/bankparse/config.toml)")

	cmd.AddCommand(NewParseCmd())
	return cmd
}

func main() {
	if err := NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
