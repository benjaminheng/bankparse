package main

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// Row represents the parsed transaction data. This is the common
// representation that different parsers will output. It is also modelled after
// YNAB's import format. The ubiquity of YNAB means that this format should
// also be supported by the majority of other budgeting tools.
type Row struct {
	Date   time.Time
	Payee  string
	Memo   string
	Amount float64 // Negative if outflow; positive if inflow
}

// Parser describes the interface that parsers should implement. A parser takes
// an io.Reader as input and returns a slice of Row structs.
type Parser interface {
	Parse(reader io.Reader) ([]Row, error)
}

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
		RunE:  parse(DBSCSVParser{}),
		Args:  cobra.ExactArgs(1),
	}

	parseDBSRawTableCmd := &cobra.Command{
		Use:   "dbs-raw-table <file>",
		Short: "DBS transaction history manually copied to the clipboard",
		Long: `DBS credit card statements aren't available for
		download as a CSV. The only way to export it is by copying the
		contents of the HTML table. This command supports parsing the
		copied contents.`,
		RunE: parse(DBSRawTableParser{}),
		Args: cobra.ExactArgs(1),
	}

	parseCmd.AddCommand(parseDBSDebitFormatCmd)
	parseCmd.AddCommand(parseDBSRawTableCmd)
	return parseCmd
}

func parse(parser Parser) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		f, err := os.Open(args[0])
		if err != nil {
			return errors.Wrap(err, "open file")
		}
		defer f.Close()

		rows, err := parser.Parse(f)
		if err != nil {
			return errors.Wrap(err, "parse DBS debit format")
		}

		w := csv.NewWriter(os.Stdout)
		w.Write([]string{"date", "payee", "memo", "amount"})
		for _, v := range rows {
			record := []string{v.Date.Format("2006-01-02"), v.Payee, v.Memo, strconv.FormatFloat(v.Amount, 'f', 2, 64)}
			w.Write(record)
		}
		w.Flush()

		return nil
	}
}
