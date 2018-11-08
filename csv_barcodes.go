package csv_barcode

import (
	"bufio"
	"encoding/csv"
	"fmt"
	barcode2 "github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code39"
	"github.com/jung-kurt/gofpdf"
	"image/jpeg"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

type Student struct {
	first string `json: first`
	last string `json: last`
	pin string	`json: pin`
}

func getPin(password string) string {
	pattern, _ := regexp.Compile("([0-9]+)")
	pinStr := pattern.FindString(password)
	return pinStr
}
//Reads CSV file. The file must comply with RFC 4180 "Common Format and MIME Type for CSV Files".
func ReadCsv(filePath string) []Student{
	var students []Student
	csvFile, _ := os.Open(filePath)
	reader := csv.NewReader(bufio.NewReader(csvFile))

	for {
		line, err := reader.Read()

		if err == io.EOF {
			break
		}else if err != nil{
			log.Fatal(err)
		}

		students = append(students, Student{
			last: line[0],
			first: strings.TrimSpace(line[1]),
			pin: getPin(line[2]),
		})
	}
	return students
}

//Makes individual png barcodes
func makeBarcodeFile(location, filename, code string)  {
	err := os.Chdir(location)
	checkError(err, "Can't change directory for barcode.")

	barcode, err := code39.Encode(code, false, true)
	checkError(err, "Can't generate barcode.")

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
		makeBarcodeFile(fileDir, filename, student.pin)
	}
}

func GeneratePdf(path, filename string, students []Student){
	MakeBarcodes(path + "/barcodes", students)
	pdf := gofpdf.New("P", "in", "Letter", "")
	pdf.SetAutoPageBreak(true, -1)
	pdf.SetMargins(0.393750, .2, 0.393750)
	left, top, right, bottom := pdf.GetMargins()
	fmt.Printf("Left: %f, Top: %f, Right: %f, Bottom: %f", left, top, right, bottom)
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 12)
	pdf.Ln(1)

	os.Chdir(path + "/barcodes")
	for _, student := range students{
		name := student.last + ", " + student.first
		barcodeName := student.last + "_" + student.first + ".jpg"
		pdf.WriteAligned(10.25,2,name, "C")
		pdf.ImageOptions(barcodeName, 4.5, 0,2,0, true,  gofpdf.ImageOptions{ImageType: "JPG"}, 0, "")
		pdf.Ln(1)
	}
	err := pdf.OutputFileAndClose(path + filename)

	if err != nil {
		log.Fatal(err.Error())
	}
}

func checkError(err error, msg string){
	if err != nil{
		fmt.Println(msg)
		log.Fatal(err)
	}
}