package csv_barcode

import (
	"bufio"
	"encoding/csv"
	"fmt"
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
//Reads CSV file. Note: each record must terminate with \n.
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

func WriteCsv(location, filename string, students []Student){
	err := os.Chdir(location)
	checkError(err, "Can't change directory.")

	file, err:=  os.Create(filename)
	checkError(err, "Cannot create file.")
	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, student := range students{
		studentInfo := []string{student.last, student.first, student.pin}
		fmt.Println(studentInfo)
		writer.Write(studentInfo)
	}
}

func checkError(err error, msg string){
	if err != nil{
		fmt.Println(msg)
		log.Fatal(err)
	}
}