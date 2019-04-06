# barcode_pdf_generator
A small library that generates a PDF of barcodes. 

## Use

### Setup of CSV File
The generator requires a csv file that must be formatted according to [RFC 4180](https://tools.ietf.org/html/rfc4180). To make sure everything works properly be sure:
	
* Values are separated by commas. 
* All records end with a CRLF line break. 
* There are **no headings** on the fields (I hope the change this in the future).

Values should be:

1. Last Name
2. First Name
3. Pin or the numeric value of the barcode. If any alphanumeric values are enterred, only numerals will be used.

### Generate Barcodes from a CSV File

1. Use the `ReadCsv(*filepath*)` method with the absolute path to the csv file to return an array of `Students` that will be needed for step 2.
2. Run `GeneratePdf(*pdfPath*, *filename*, *[]student*)` using the full destination path, filename, and the `[]Students` from step 1. 

**At this time, no heading can be added to the PDF.**


