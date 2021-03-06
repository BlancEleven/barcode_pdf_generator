package csv_barcode

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	barcode2 "github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code39"
	"github.com/jung-kurt/gofpdf"
)

type Student struct {
	first  string
	last   string
	pinSet bool
	pin    string
}

func getPin(password string) (string, bool) {
	pattern, _ := regexp.Compile("([0-9]+)")
	pinStr := pattern.FindString(password)
	//If pin isn't set.
	if pinStr == "" {

		return "", false
	}
	return pinStr, true
}

//Reads CSV file. The file must comply with RFC 4180 "Common Format and MIME Type for CSV Files".
//This function will return []Student that's needed to run GeneratePdf
func ReadCsv(filePath string) []Student {
	var students []Student
	csvFile, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening CSV file: \n%s", err)
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))

	for {
		line, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		pin, pinSet := getPin(line[2])

		students = append(students, Student{
			last:   line[0],
			first:  strings.TrimSpace(line[1]),
			pin:    pin,
			pinSet: pinSet,
		})
	}
	return students
}
func dirExists(location string) bool {
	if _, err := os.Stat(location); err != nil {
		return false
	}
	return false
}

//Makes individual png barcodes
func makeBarcodeFile(location, filename, code string, pinSet bool) {
	var barcode barcode2.Barcode
	err := os.Chdir(location)
	checkError(err, "Can't change directory for individual barcodes.")

	if pinSet {
		barcode, err = code39.Encode(code, false, true)
		checkError(err, "Can't generate barcode.")

	} else {
		barcode, err = code39.Encode("NONE", false, true)
		checkError(err, "Can't generate barcode.")
	}

	scaled, err := barcode2.Scale(barcode, 250, 100)
	checkError(err, "Error scaling barcode.")

	file, err := os.Create(filename)
	checkError(err, "Cannot create barcode file.")

	defer file.Close()
	jpeg.Encode(file, scaled, nil)

}

//Generates Barcodes to the requested directory
func MakeBarcodes(fileDir string, records []Student) {
	for _, student := range records {

		filename := student.last + "_" + student.first + ".jpg"
		makeBarcodeFile(fileDir, filename, student.pin, student.pinSet)
	}
}

//Generates a pdf with barcodes and names below them.
func GeneratePdf(pdfPath, filename, heading string, students []Student) {
	if !dirExists(pdfPath) {
		os.Mkdir(pdfPath, 0700)
		os.Mkdir(pdfPath+"/barcodes", 0700)
	}

	MakeBarcodes(pdfPath+"/barcodes", students)
	pdf := gofpdf.New("P", "in", "Letter", "")
	pdf.SetAutoPageBreak(true, -1)
	pdf.SetMargins(-0.75, .2, 0)
	pdf.AddPage()
	//Heading
	pdf.SetFont("Arial", "B", 25)
	pdf.WriteAligned(10.25, 2, heading, "C")
	pdf.Ln(1)
	//
	pdf.SetFont("Arial", "B", 12)
	pdf.Ln(1)

	os.Chdir(pdfPath + "/barcodes")

	for _, student := range students {
		var name string
		if student.pinSet {
			name = student.first + " " + student.last
		} else {
			name = student.first + " " + student.last + " - NO PIN SET"
		}

		barcodeName := student.last + "_" + student.first + ".jpg"
		pdf.WriteAligned(10.25, 2, name, "C")
		pdf.ImageOptions(barcodeName, 3.35, 0, 2, 0, true, gofpdf.ImageOptions{ImageType: "JPG"}, 0, "")
		pdf.Ln(1)
	}
	err := pdf.OutputFileAndClose("../" + filename)

	if err != nil {
		log.Fatal(err.Error())
	}
}

func checkError(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		log.Fatal(err)
	}
}
