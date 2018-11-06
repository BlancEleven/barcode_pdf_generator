package csv_barcode

import (
	"bufio"
	"encoding/csv"
	"fmt"
	barcode2 "github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code39"
	"image/png"
	"io"
	"log"
	"os"
	"regexp"
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
			first: line[1],
			pin: getPin(line[2]),
		})
	}
	return students
}

//Makes individual png barcodes
func MakeBarcodeFile(location, filename string, code string)  {
	err := os.Chdir(location)
	checkError(err, "Can't change directory for barcode.")
	barcode, err := code39.Encode(code, false, true)
	checkError(err, "Can't generate barcode.")
	scaled, err := barcode2.Scale(barcode, 250, 100)
	checkError(err, "Error scaling barcode.")
	file, err := os.Create(filename)
	checkError(err, "Cannot create barcode file.")
	defer file.Close()
	png.Encode(file, scaled)

}

//Generates Barcodes to the requested directory
func MakeBarcodes(fileDir string, records []Student) {
	for _, student := range records {
		filename := student.last + "_" + student.first + ".png"
		MakeBarcodeFile(fileDir, filename, student.pin)
	}
}


func checkError(err error, msg string){
	if err != nil{
		fmt.Println(msg)
		log.Fatal(err)
	}
}