# csvtx

Import & export financial transactions in csv format.

Use case: exporting from Mint & importing to YNAB4 (classic edition) is a pain
because the fields don't match up too well. This program takes a csv file as
input and outputs a formatted csv file suited for import into YNAB4.

You *could* use a spreadsheet to transform your data, or you can do it with code!

### usage

From you Mint account, [export
transactions](https://help.mint.com/Accounts-and-Transactions/888960591/How-do-you-export-transaction-data.htm)
to CSV. Pass in the file path of the exported csv file your function. Run the
program, which produces another csv file for every unique account name found
among the rows in input csv file. Use each output file to populate the
corresponding account in YNAB4.
