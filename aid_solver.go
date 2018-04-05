package main

import (
	"fmt"
	"flag"
	"encoding/csv"
	"os"
	"log"
	"bufio"
	"io"
	"regexp"
)

const phpKeyRegexp = `[\'|\"](\S+)[\'|\"][\s+]?\=\>[\s+]?[\'|\"](\d+)[\'|\"]`

func main() {
	csvFilePath := flag.String("csvPath", "", "Path to csv file")
	phpFilePath := flag.String("phpPath", "", "Path to php file")
	csvEnVariablesPosition := flag.Int("enVarsPos", 5, "Position of en variables in csv file, starts from 0")
	csvCustomVariablesPosition := flag.Int("customVarsPos", 6, "Position of custom variables in csv file, starts from 0")

	flag.Parse()

	if _, err := os.Stat(*csvFilePath); os.IsNotExist(err) {
		fmt.Printf("File \"%s\" is not found\n", *csvFilePath)
		return
		os.Exit(0)
	}

	fmt.Printf("Read csv file...\n")

	csvDictionary := readCsvFile(*csvFilePath, *csvEnVariablesPosition, *csvCustomVariablesPosition)
	phpDictionary, _ := readPHPArray(*phpFilePath)

	for originalAidValue, parameter := range phpDictionary {
		for enValue, customValue := range csvDictionary {
			if originalAidValue == enValue {
				fmt.Printf("'%s' => '%s',\n ", parameter, customValue);
			}
		}
	}
}

func readCsvFile(path string, enPosition int, customPosition int) map[string]string {
	file := openFile(path)
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	resultMap := make(map[string]string)

	for _, v := range records {
		enValue := ""
		customValue := ""
		for key, value := range v {
			if key == enPosition && len(value) > 1 {
				enValue = value
			}
			if key == customPosition && len(value) > 1 {
				customValue = value
			}
		}

		if customValue != "" && enValue != "" {
			resultMap[enValue] = customValue
			//fmt.Println("en:", enValue, "custom:", customValue)
		}
	}
	return resultMap
}

func openFile(path string) (*os.File) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}


	return file
}

func readPHPArray(path string) (map[string]string, error) {
	f := openFile(path)

	defer f.Close()

	r := bufio.NewReader(f)

	resultMap := make(map[string]string)

	var regExp = regexp.MustCompile(phpKeyRegexp)

	for {
		lineString, err := r.ReadString(10) // 0x0A separator = newline
		if err == io.EOF {
			break
		} else if err != nil {
			return resultMap, err // if you return error
		}

		aidKey := ""
		aidValue := ""
		for i, match := range regExp.FindStringSubmatch(lineString) {
			//fmt.Println(match, "found at index", i)
			if i == 1 && len(match) > 1 {
				aidKey = match
			}
			if i == 2 && len(match) > 1 {
				aidValue = match
			}
		}



		if aidKey != "" && aidValue != "" {
			resultMap[aidValue] = aidKey;
			//fmt.Println("key:", aidKey, " ===> value:", aidValue)
		}
	}

	return resultMap, nil
}