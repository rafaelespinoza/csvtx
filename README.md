# csvtx

Convert exported CSV financial transactions for import into YNAB4.

Currently supported source services:

- Mint

Motivation: I like YNAB4 (classic edition) for personal budgeting, but it
requires you to manually input each transaction and I don't always keep up with
that. I wanted a tool to take some exported CSV data from a financial
institution, or a service such as Mint, and format it so I can do a bulk import
into YNAB4.

### usage

Export a CSV of financial transactions from a source service (supported ones
listed above). Save the CSV somewhere.

Get dependencies and build a binary.
```sh
make all
```

Then use the `convert` command to convert a CSV.

Use `convert mint-to-ynab` to convert exported transactions from Mint.com. It
will create one output file per account type.

```sh
% ./bin/csvtx convert mint-to-ynab -i path/to/transactions.csv
wrote "Checking" file "/tmp/checking.csv"
wrote "PERSONAL SAVINGS" file "/tmp/personal-savings.csv"
wrote "Credit" file "/tmp/credit.csv"
```

The `-o` flag can set the directory of the output file.
