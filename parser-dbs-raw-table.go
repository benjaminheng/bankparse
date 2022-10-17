package main

import "io"

// DBSRawTableParser parses a HTML table copied to clipboard. This will
// typically be similar to a tab-separated CSV. This is used for credit card
// statements, which DBS does not provide a proper CSV export for.
type DBSRawTableParser struct{}

func (f DBSRawTableParser) Parse(reader io.Reader) ([]Row, error) {
	return nil, nil
}
