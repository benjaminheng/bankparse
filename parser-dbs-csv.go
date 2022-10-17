package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// DBSCSVParser parses the exported DBS CSV file.
type DBSCSVParser struct{}

func (f DBSCSVParser) parseRow(record []string, loc *time.Location) (Row, error) {
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

func (f DBSCSVParser) Parse(reader io.Reader) ([]Row, error) {
	var csvContents string
	var csvHeaderSeen bool

	scanner := bufio.NewScanner(reader)
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
		return nil, err
	}
	csvContents = strings.TrimSpace(csvContents)

	r := csv.NewReader(strings.NewReader(string(csvContents)))
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
