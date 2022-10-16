package main

import "github.com/spf13/cobra"

func NewParseCmd() *cobra.Command {
	parseCmd := &cobra.Command{
		Use:   "parse",
		Short: "Parses a bank's transaction history",
		Long:  ``,
	}

	parseDBSDebitFormatCmd := &cobra.Command{
		Use:   "dbs-csv",
		Short: "DBS transaction history exported as CSV",
		Long:  `DBS debit card transaction history can be exported as CSV.`,
		RunE:  parseDBSDebitFormat,
	}

	parseDBSRawTableCmd := &cobra.Command{
		Use:   "dbs-raw-table",
		Short: "DBS transaction history manually copied to the clipboard",
		Long: `DBS credit card statements aren't available for
		download as a CSV. The only way to export it is by copying the
		contents of the HTML table. This command supports parsing the
		copied contents.`,
		RunE: parseDBSRawTableFormat,
	}

	parseCmd.AddCommand(parseDBSDebitFormatCmd)
	parseCmd.AddCommand(parseDBSRawTableCmd)
	return parseCmd
}

func parseDBSDebitFormat(cmd *cobra.Command, args []string) error {
	return nil
}

func parseDBSRawTableFormat(cmd *cobra.Command, args []string) error {
	return nil
}
