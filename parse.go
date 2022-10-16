package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewParseCmd() *cobra.Command {
	parseCmd := &cobra.Command{
		Use:   "parse",
		Short: "Parses a bank's transaction history",
		Long:  ``,
	}

	parseDBSDebitFormatCmd := &cobra.Command{
		Use:   "dbs-csv <file>",
		Short: "DBS transaction history exported as CSV",
		Long:  `DBS debit card transaction history can be exported as CSV.`,
		RunE:  parseDBSDebitFormat,
		Args:  cobra.ExactArgs(1),
	}

	parseDBSRawTableCmd := &cobra.Command{
		Use:   "dbs-raw-table <file>",
		Short: "DBS transaction history manually copied to the clipboard",
		Long: `DBS credit card statements aren't available for
		download as a CSV. The only way to export it is by copying the
		contents of the HTML table. This command supports parsing the
		copied contents.`,
		RunE: parseDBSRawTableFormat,
		Args: cobra.ExactArgs(1),
	}

	parseCmd.AddCommand(parseDBSDebitFormatCmd)
	parseCmd.AddCommand(parseDBSRawTableCmd)
	return parseCmd
}

func parseDBSDebitFormat(cmd *cobra.Command, args []string) error {
	f, err := os.Open(args[0])
	if err != nil {
		return errors.Wrap(err, "open file")
	}
	defer f.Close()

	var csvContents string
	var csvHeaderSeen bool

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Transaction Date,") {
			csvHeaderSeen = true
			csvContents += line + "\n"
		} else if csvHeaderSeen && line != "" {
			csvContents += line + "\n"
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	csvContents = strings.TrimSpace(csvContents)
	fmt.Println(csvContents)
	return nil
}

func parseDBSRawTableFormat(cmd *cobra.Command, args []string) error {
	return nil
}
