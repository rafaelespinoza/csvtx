# csvtx

[![ci](https://github.com/rafaelespinoza/csvtx/actions/workflows/ci.yml/badge.svg)](https://github.com/rafaelespinoza/csvtx/actions/workflows/ci.yml)

Convert exported CSV financial transactions for import into YNAB4.

Currently supported source services:

- Mechanic's Bank
- Mint
- Wells Fargo
- Venmo

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

Convert exported Mint.com transactions and produce one output CSV file per
account type.

```sh
% ./bin/csvtx convert -from mint -i path/to/transactions.csv
wrote "Checking" file "/tmp/checking.csv"
wrote "PERSONAL SAVINGS" file "/tmp/personal-savings.csv"
wrote "Credit" file "/tmp/credit.csv"
```

Convert Mechanic's Bank transactions.
```sh
% ./bin/csvtx convert -from mechanicsbank -i path/to/transactions.csv
wrote "mechanicsbank" file "/tmp/mechanicsbank.csv"
```

Convert Wells Fargo transactions.
```sh
% ./bin/csvtx convert -from wellsfargo -i path/to/transactions.csv
wrote "wellsfargo" file "/tmp/wellsfargo.csv"
```

Convert Venmo transactions.
```sh
% ./bin/csvtx convert -from venmo -i path/to/transactions.csv
wrote "venmo" file "/tmp/venmo.csv"
```

The `-o` flag can set the directory of the output file.
