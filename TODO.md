## TODO

- [ ] Support `--config` flag
- [ ] Add `init` command. Write config to `$XDG_CONFIG_HOME/bankparse/config.toml` or `$HOME/.config/bankparse/config.toml`
- [ ] Implement parser for DBS csv format
- [ ] Implement parser for DBS raw table format
- [ ] Write parsed output to csv file. Add idempotency.
- [ ] Implement categorization logic as new `categorize` command.
    - [ ] Add `categorization_attempted` column so we can identify and
      reconcile the transactions that couldn't be categorized.
    - [ ] Add `--rerun` flag to rerun categorization logic on already
      categorized columns. We can use this flag after updating rules to cover
      previously uncovered transactions.
- [ ] Add `--categorize` flag to `parse` command. Invoke categorization logic when enabled.

## Open questions

- How to handle cash transactions? Maybe we can just mark ATM withdrawals as
  such, and ignore them in our budgeting since we don't use cash enough for it
  to matter.
