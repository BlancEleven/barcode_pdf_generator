package csv_barcode

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"

	//"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type Student struct {
	first string
	last string
	pin int
}

func getPin(password string) int {
	pattern, _ := regexp.Compile("([0-9]+)")

	pinStr := pattern.FindString(password)
	pin, err := strconv.Atoi(pinStr)
		if err != nil {
			fmt.Printf("There was an error extracting the pin: %s", err)
		}
	return pin
}

func ReadFile(filePath string)  {
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
	fmt.Println(students)
}