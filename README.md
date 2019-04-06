# barcode_pdf_generator
A small library that generates a PDF of barcodes. 

## Use

###Setup of CSV File
The generator requires a csv file that must be formatted according to [RFC 4180] (https://tools.ietf.org/html/rfc4180). To make sure everything works properly be sure:
	
	* Values are separated by commas. 
	* All records must end with a CRLF line break. 
	* There are **no headings** on the fields (I hope the change this in the future). 


