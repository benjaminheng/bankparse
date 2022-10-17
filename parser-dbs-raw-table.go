package main

import (
	"encoding/csv"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// DBSRawTableParser parses a HTML table copied to clipboard. This will
// typically be similar to a TSV file. This parser can be used for credit card
// statements, which DBS does not provide a proper CSV export for.
type DBSRawTableParser struct{}

func (f DBSRawTableParser) parseRow(record []string, loc *time.Location) (Row, error) {
	rawDate := strings.TrimSpace(record[0])
	payee := strings.TrimSpace(record[1])
	rawAmount := strings.TrimSpace(record[2])

	// Parse date
	date, err := time.ParseInLocation("02 Jan 2006", rawDate, loc)
	if err != nil {
		return Row{}, errors.Wrap(err, "parse date")
	}

	// Strip the currency symbol
	if components := strings.Split(rawAmount, "$"); len(components) > 1 {
		rawAmount = components[1]
	}

	// Parse amount to float
	amount, err := strconv.ParseFloat(rawAmount, 64)
	if err != nil {
		return Row{}, errors.Wrap(err, "parse amount to float64")
	}
	amount = -amount // Credit card transactions are always outflow

	row := Row{
		Date:   date,
		Payee:  payee,
		Amount: -amount,
	}
	return row, nil
}

func (f DBSRawTableParser) Parse(reader io.Reader) ([]Row, error) {
	r := csv.NewReader(reader)
	r.Comma = '\t'
	records, err := r.ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, "read csv contents")
	}

	loc, err := time.LoadLocation("Asia/Singapore")
	if err != nil {
		return nil, errors.Wrap(err, "load local timezone")
	}

	var rows []Row
	for _, record := range records {
		row, err := f.parseRow(record, loc)
		if err != nil {
			return nil, errors.Wrap(err, "parse row")
		}
		rows = append(rows, row)
	}
	return rows, nil
}
