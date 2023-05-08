package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gocarina/gocsv"
)

// Put cfg data with the struct? I guess?
var BOX_LEN = 7 // work around for trailing whitespace...
var OFFSET_INDEXES = []int{1, 4}

type Box struct {
	User_id  string `csv:"user_id"`
	Num_foo  int32  `csv:"num_foo"`
	Foos     Foober `csv:"foos"`
	Mew      string `csv:"mew"`
	Num_bars int8   `csv:"Num_bars"`
	Bars     Foober `csv:"bars"`
	Mow      string `csv:"Mow!"`
}

// aslkdfjakl;jds
type Foober struct {
	foos []string
}

// couldn't figure out how to unmarshal directly into []string
func (foo *Foober) UnmarshalCSV(csv string) (err error) {
	foo.foos = strings.Split(csv, " ")
	return err
}

func (foo Foober) MarshalJSON() ([]byte, error) {
	return json.Marshal(foo.foos)
}

func main() {
	file, err := os.Open("data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Wanted this out of the loop but it was easier to read
	// when it was down where it's used
	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.Comma = '😀'
		return r
	})

	scanner := bufio.NewScanner(file)
	// We read file line by line
	for scanner.Scan() {
		// We need a new scanner every time (or maybe we can reset it idk)
		r := csv.NewReader(strings.NewReader(scanner.Text()))
		r.Comma = ' '
		row, err := r.Read()
		if err != nil {
			fmt.Println("Error parsing line as (csv) row:", err)
		}

		// Preprocess line for each offset (or, the indexes of each num_foo field)
		for _, offset := range OFFSET_INDEXES {
			num_fields, err := strconv.Atoi(row[offset])
			if err != nil {
				fmt.Println("error: ", err)
				return
			}
			// It will chop off data when we get a 0 without checking this
			if num_fields != 0 {
				offset += 1
			}
			squashed_fields := strings.Join(row[offset:num_fields+offset], " ")
			temp := append(row[:offset], squashed_fields)
			// Modify row for next pass or step
			row = append(temp, row[num_fields+offset:]...)
		}
		// Preprocessing done
		normalized_string := strings.Join(row[:BOX_LEN], "😀")
		storage := [1]Box{} // Only a single row since we're going line by line
		reader := strings.NewReader(normalized_string)
		err = gocsv.UnmarshalWithoutHeaders(reader, &storage)
		if err != nil {
			fmt.Println("Err,line:", err, normalized_string)
		}
		output, err := json.Marshal(storage)
		fmt.Printf("%+v\n", string(output)) // Line is done
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
