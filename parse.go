package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type Row struct {
	Date   time.Time
	Payee  string
	Memo   string
	Amount float64
}

type DBSDebitFormat string

func (f DBSDebitFormat) parseRow(record []string, loc *time.Location) (Row, error) {
	// Headers: Transaction Date,Reference,Debit Amount,Credit Amount,Transaction Ref1,Transaction Ref2,Transaction Ref3
	rawDate := strings.TrimSpace(record[0])
	debitAmount := strings.TrimSpace(record[2])
	creditAmount := strings.TrimSpace(record[3])
	payee := strings.TrimSpace(record[4])

	// Parse date
	date, err := time.ParseInLocation("02 Jan 2006", rawDate, loc)
	if err != nil {
		return Row{}, errors.Wrap(err, "parse date")
	}

	// Construct memo by concatenating the transaction ref2 and ref3
	var memo string
	var memos []string
	if record[5] != "" {
		memos = append(memos, record[5])
	}
	if record[6] != "" {
		memos = append(memos, record[6])
	}
	if len(memos) > 0 {
		memo = strings.Join(memos, ",")
	}

	// Construct amount from debit and credit
	var amount float64
	if debitAmount != "" {
		debit, err := strconv.ParseFloat(debitAmount, 64)
		if err != nil {
			return Row{}, errors.Wrap(err, "parse debit amount to float64")
		}
		amount = -debit
	} else if creditAmount != "" {
		credit, err := strconv.ParseFloat(creditAmount, 64)
		if err != nil {
			return Row{}, errors.Wrap(err, "parse credit amount to float64")
		}
		amount = credit
	} else {
		return Row{}, errors.New("neither debit nor credit amount set")
	}

	row := Row{
		Date:   date,
		Amount: amount,
		Payee:  payee,
		Memo:   memo,
	}
	return row, nil
}

func (f DBSDebitFormat) Parse() ([]Row, error) {
	r := csv.NewReader(strings.NewReader(string(f)))
	records, err := r.ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, "read csv contents")
	}

	loc, err := time.LoadLocation("Asia/Singapore")
	if err != nil {
		return nil, errors.Wrap(err, "load local timezone")
	}

	var rows []Row
	for i, record := range records {
		if i == 0 {
			continue // skip header
		}
		row, err := f.parseRow(record, loc)
		if err != nil {
			return nil, errors.Wrap(err, "parse row")
		}
		rows = append(rows, row)
	}
	return rows, nil
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
			// DBS appends an empty value to each record. Header
			// has 7 elements, record has 8 elements. Remove the
			// extra element from each record.
			if line[len(line)-1] == ',' {
				line = line[:len(line)-1]
			}
			csvContents += line + "\n"
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	csvContents = strings.TrimSpace(csvContents)

	dbsDebitFormat := DBSDebitFormat(csvContents)
	rows, err := dbsDebitFormat.Parse()
	if err != nil {
		return errors.Wrap(err, "parse DBS debit format")
	}
	for _, v := range rows {
		fmt.Println(v)
	}

	return nil
}

func parseDBSRawTableFormat(cmd *cobra.Command, args []string) error {
	return nil
}
