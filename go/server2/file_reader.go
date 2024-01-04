package main

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
	"github.com/cathieyun/libfss/libfss"
)

// Runs FSS to evalute db records
func readIntFetchSmall(server *libfss.Fss, serverNum byte, fssKey libfss.FssKeyEq2P, fileName string) string {
	var ans int = 0
	file, err := os.Open("tpch_line_number.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return ""
	}
	defer file.Close()

	reader := csv.NewReader(file)
	
	// Read the header (first line) and discard it
	_, err = reader.Read()
	if err != nil {
		fmt.Println("Error reading header:", err)
		return ""
	}


	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return ""
	}
	startTime := time.Now()
	// Process each line
	for _, record := range records {
		// print(i)
		uuint, err := strconv.ParseUint(record[1], 10, 0)
		print(uuint)
		if err != nil {
			return ""
		}

		ans += server.EvaluatePF(serverNum, fssKey, uint(uuint))

	}
	endTime := time.Now()
	print("evaluation time  for 1m-row in nanosecond :", endTime.Sub(startTime).Nanoseconds())

	return strconv.Itoa(ans)

}

func indexOf(column string, header []string) int {
	for i, col := range header {
		if col == column {
			return i
		}
	}
	return -1
}


func stringToInt(s string) uint {
	h := sha256.New()
	h.Write([]byte(s))
	num := binary.LittleEndian.Uint32(h.Sum(nil))
	return uint(num)
}
