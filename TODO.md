## TODO

- [x] Implement parser for DBS csv format
- [ ] Implement parser for DBS raw table format
- [x] Write parsed output to stdout in csv format
- [ ] Add `parse -i` flag for interactive input, will probably be used for the
  DBS raw table format.

## Open questions

- How to handle cash transactions? Maybe we can just mark ATM withdrawals as
  such, and ignore them in our budgeting since we don't use cash enough for it
  to matter.
- Do we need to support visualizations? Such as showing aggregations and
  summaries. Might be easier to just defer that functionality to Google Sheets
  or VisiData.
- Should we handle automatic categorizations? Right now I'm leaning towards no.
  I'll probably be using ActualBudget for my budgeting, and I just need
  something to format DBS transactions into a usable CSV format.
